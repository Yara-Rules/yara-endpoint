package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/common/errors"
	"github.com/Yara-Rules/yara-endpoint/common/message"
	"github.com/Yara-Rules/yara-endpoint/server/config"
	"github.com/Yara-Rules/yara-endpoint/server/database"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	"github.com/k0kubun/pp"
	"github.com/oklog/ulid"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	conn net.Conn
	w    *gob.Encoder
	r    *gob.Decoder
	L    net.Listener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Connect(c *config.Config) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.TCPServer.Addr, c.TCPServer.Port))
	checkErr(err)

	s.L, err = net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
}

func (s *Server) Close() {
	s.conn.Close()
}

func (s *Server) Accept() (net.Conn, error) {
	var err error
	s.conn, err = s.L.Accept()
	if err != nil {
		return s.conn, err
	}
	s.r = gob.NewDecoder(s.conn)
	s.w = gob.NewEncoder(s.conn)
	return s.conn, err
}

func (s *Server) HandleClient() {
	defer s.Close()

	db := DB.NewDataStore()
	defer db.Close()

	var msg message.Message
	s.r.Decode(&msg)

	if msg.CMD == command.Reserved {
		return
	}

	switch msg.CMD {
	case command.Register:
		if msg.ULID != "" {
			log.Infof("[%s] Processing Register command", msg.ULID)
		} else {
			log.Infof("[%s] Processing Register command", s.conn.RemoteAddr())
		}
		s.processRegister(msg, db)
	case command.Ping:
		if msg.ULID != "" {
			log.Infof("[%s] Processing Ping command", msg.ULID)
		} else {
			log.Infof("[%s] Processing Ping command", s.conn.RemoteAddr())
		}
		s.processPing(msg, db)

	case command.ScanFile:
		if msg.ULID != "" {
			log.Infof("[%s] Processing ScanFile command", msg.ULID)
		} else {
			log.Infof("[%s] Processing ScanFile command", s.conn.RemoteAddr())
		}
		s.processScanFile(msg, db)

	case command.ScanDir:
		if msg.ULID != "" {
			log.Infof("[%s] Processing ScanDir command", msg.ULID)
		} else {
			log.Infof("[%s] Processing ScanDir command", s.conn.RemoteAddr())
		}
		s.processScanDir(msg, db)

	case command.ScanPID:
		if msg.ULID != "" {
			log.Infof("[%s] Processing ScanPID command", msg.ULID)
		} else {
			log.Infof("[%s] Processing ScanPID command", s.conn.RemoteAddr())
		}
		s.processScanPID(msg, db)
	}
}

func (s *Server) processRegister(msg message.Message, db *database.DataStore) {
	log.Debugf("[%s] Processing register command", msg.ULID)

	// TODO: Check Hostname value

	// New endpoint
	id := newULID()
	ep := &models.Endpoint{
		ULID:          id.String(),
		Hostname:      strings.Split(msg.Data, "|")[0],
		ClientVersion: strings.Split(msg.Data, "|")[1],
		Tags:          []string{},
		LastPing:      time.Now(),
		CreateAt:      time.Now(),
		UpdateAt:      time.Now(),
	}
	// TODO: Check potential errors
	db.C(models.Endpoints).Insert(ep)
	msg.ULID = id.String()
	s.w.Encode(msg)
}

func (s *Server) processPing(msg message.Message, db *database.DataStore) {
	log.Debugf("[%s] Processing ping command", msg.ULID)

	err := db.C(models.Endpoints).Update(bson.M{"ulid": msg.ULID}, bson.M{"$set": bson.M{"last_ping": time.Now()}})
	if err == mgo.ErrNotFound {
		log.Errorf("[%s]  Sending PONG with error: %s", msg.ULID, errors.Errors[errors.NeedsRegister])
		s.sendErrorMsg(msg, errors.NeedsRegister, "Pong")
	} else {
		if jobs, pending := s.pendingJobs(msg, db); pending {
			s.taskPicker(msg, jobs, db)
		} else {
			msg.Result = "Pong"
			s.w.Encode(msg)
			log.Infof("[%s] Sending PONG due to no task assigned", msg.ULID)
		}
	}
}

func (s *Server) taskPicker(msg message.Message, jobs models.Schedule, db *database.DataStore) {
	log.Debugf("[%s] Picking up a task ", msg.ULID)
	for _, job := range jobs.Tasks {
		// Choose the first unprocessed task which is scheduled to be executed before now
		if job.Status == models.Initial && time.Now().After(job.When) {
			msg.TaskID = job.TaskID
			msg.CMD = job.Command
			msg.Params = job.Params
			log.Debugf("[%s] Extracting rules for taskID %s", msg.ULID, msg.TaskID)
			msg.Data = s.extractRules(job.Rules, db)

			// Update Task to running
			err := s.updateTaskStatus(msg, models.Running, db)
			// selector := bson.M{"$and": []bson.M{
			// 	bson.M{"ulid": msg.ULID},
			// 	bson.M{"tasks": bson.M{"$elemMatch": bson.M{"task_id": msg.TaskID}}}}}
			// update := bson.M{"$set": bson.M{"tasks.$.status": models.Running}}

			// err := db.C(models.Schedules).Update(selector, update)
			if err == mgo.ErrCursor || err == mgo.ErrNotFound {
				// Track error
				log.Errorf("[%s] %s", msg.ULID, errors.Errors[errors.UnableToUpdateDB])
				s.saveErr(msg.ULID, msg.TaskID, errors.UnableToUpdateDB, db)
				// Send pong as usual. For now
				msg.Result = "Pong"
				s.w.Encode(msg)
				log.Infof("[%s] Sending PONG due to error while updating db", msg.ULID)
				// Quit
				return
			}
			// Send task to endpoint
			log.Infof("[%s] Selected task %s", msg.ULID, msg.TaskID)
			pp.Println(msg)
			s.w.Encode(msg)
			return
		}
	}
	log.Debugf("[%s] No tasks to be picked up", msg.ULID)
	msg.Result = "Pong"
	s.w.Encode(msg)
	log.Infof("[%s] Sending PONG due to no task assigned", msg.ULID)
	return
}

func (s *Server) processScanFile(msg message.Message, db *database.DataStore) {
	log.Debugf("[%s] Processing ScanFile response", msg.ULID)
	if s.checkMsgErr(msg, models.Failed, db) {
		// TODO: Can I recover this? Should I ignore?
		msg.Result = "KO"
	} else {
		if s.saveResult(msg, db) {
			msg.Result = "OK"
		} else {
			// TODO: Request a new scan due errors while saving the report
			log.Errorf("[%s] Error on processScanFile", msg.ULID)
		}
	}
	// TODO:
	// - check error
	// - save the result
	// - send a message back
	s.w.Encode(msg)
}

func (s *Server) processScanDir(msg message.Message, db *database.DataStore) {
	log.Debugf("[%s] Processing ScanDir response", msg.ULID)
	if s.checkMsgErr(msg, models.Failed, db) {
		// TODO: Can I recover this? Should I ignore?
		msg.Result = "KO"
	} else {
		if s.saveResult(msg, db) {
			msg.Result = "OK"
		} else {
			// TODO: Request a new scan due errors while saving the report
			log.Errorf("[%s] Error on processScanDir", msg.ULID)
		}
	}
	// TODO:
	// - check error
	// - save the result
	// - send a message back
	s.w.Encode(msg)
}

func (s *Server) processScanPID(msg message.Message, db *database.DataStore) {
	log.Debugf("[%s] Processing ScanPID response", msg.ULID)
	if s.checkMsgErr(msg, models.Failed, db) {
		// TODO: Can I recover this? Should I ignore?
		msg.Result = "KO"
	} else {
		if s.saveResult(msg, db) {
			msg.Result = "OK"
		} else {
			// TODO: Request a new scan due errors while saving the report
			log.Errorf("[%s] Error on processScanPID", msg.ULID)
		}
	}
	// TODO:
	// - check error
	// - save the result
	// - send a message back
	s.w.Encode(msg)
}

func (s *Server) sendErrorMsg(msg message.Message, e errors.Error, r string) {
	log.Debugf("[%s] Sending error message", msg.ULID)
	msg.Result = r
	msg.Error = true
	msg.ErrorID = e
	msg.ErrorMsg = errors.Errors[e]
	msg.TaskID = ""
	msg.Params = ""
	msg.Data = ""
	s.w.Encode(msg)
}

func (s *Server) pendingJobs(msg message.Message, db *database.DataStore) (models.Schedule, bool) {
	log.Debugf("[%s] Has it pending jobs?", msg.ULID)
	var schedule models.Schedule
	err := db.C(models.Schedules).Find(bson.M{"ulid": msg.ULID}).One(&schedule)
	if err == mgo.ErrNotFound {
		log.Debugf("[%s] No pending jobs", msg.ULID)
		return schedule, false
	}
	log.Debugf("[%s] Pending jobs", msg.ULID)
	return schedule, true
}

func (s *Server) extractRules(RuleList []bson.ObjectId, db *database.DataStore) string {
	var rules string = ""
	var rule models.Rule
	for _, ruleID := range RuleList {
		err := db.C(models.Rules).Find(bson.M{"_id": ruleID}).One(&rule)
		if err == mgo.ErrNotFound {
			log.Errorf("Rule %s not found", ruleID)
			continue
		}
		rules += rule.Data
	}
	return rules
}

func (s *Server) checkMsgErr(msg message.Message, state models.State, db *database.DataStore) bool {
	log.Debugf("[%s] Check error in message", msg.ULID)
	pp.Println(msg)
	if msg.Error {
		log.Errorf("[%s] %s", msg.ULID, errors.Errors[msg.ErrorID])
		s.saveErr(msg.ULID, msg.TaskID, msg.ErrorID, db)
		err := s.updateTaskStatus(msg, state, db)
		if err != nil {
			log.Debugf("[%s] Unable to update task %s to status %s", msg.ULID, msg.TaskID, models.States[state])
		}
		return true
	}
	return false
}

func (s *Server) saveErr(ulid, taskID string, err errors.Error, db *database.DataStore) {
	log.Debugf("[%s] Saving error for TaskID: %s", ulid, taskID)
	e := &models.Error{
		ULID:     ulid,
		TaskID:   taskID,
		ErrorID:  err,
		ErrorMsg: errors.Errors[err],
		CreateAt: time.Now(),
	}
	db.C(models.Errors).Insert(e)
}

func (s *Server) saveResult(msg message.Message, db *database.DataStore) bool {
	log.Debugf("[%s] Saving result for TaskID %s", msg.ULID, msg.TaskID)
	var r models.Report
	var results []models.Result
	err := db.C(models.Reports).Find(bson.M{"ulid": msg.ULID}).One(&r)
	if err == mgo.ErrNotFound { // First report
		r.ULID = msg.ULID
		r.CreateAt = time.Now()
		r.UpdateAt = time.Now()
		for file, matches := range msg.ResultYara {
			for _, match := range matches {
				result := models.Result{
					File:      file,
					RuleName:  match.Rule,
					Namespace: match.Namespace,
					Tags:      match.Tags,
					Meta:      match.Meta,
				}
				for _, rule := range match.Strings {
					result.Strings = append(result.Strings, models.String{
						Name:   rule.Name,
						Offset: rule.Offset,
					})
				}
				results = append(results, result)
			}
		}
		r.Reports = append([]models.LReports{}, models.LReports{
			ReportID: newULID().String(),
			Task:     s.extractTask(msg.ULID, msg.TaskID, db),
			Result:   results,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		})

		err = db.C(models.Reports).Insert(&r)
		if err != nil {
			log.Debugf("[%s] Unable to insert report for TaskID: %s", msg.ULID, msg.TaskID)
			return false
		}
	} else {
		r.UpdateAt = time.Now()
		for file, matches := range msg.ResultYara {
			for _, match := range matches {
				result := models.Result{
					File:      file,
					RuleName:  match.Rule,
					Namespace: match.Namespace,
					Tags:      match.Tags,
					Meta:      match.Meta,
				}
				for _, rule := range match.Strings {
					result.Strings = append(result.Strings, models.String{
						Name:   rule.Name,
						Offset: rule.Offset,
					})
				}
				results = append(results, result)
			}
		}
		r.Reports = append(r.Reports, models.LReports{
			ReportID: newULID().String(),
			Task:     s.extractTask(msg.ULID, msg.TaskID, db),
			Result:   results,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		})
		err = db.C(models.Reports).Update(bson.M{"ulid": r.ULID}, &r)
		if err != nil {
			log.Debugf("[%s] Unable to insert report for TaskID: %s", msg.ULID, msg.TaskID)
			return false
		}
	}
	log.Debugf("[%s] Result saved for task: %s", msg.ULID, msg.TaskID)
	return true
}

func (s *Server) extractTask(ulid, taskID string, db *database.DataStore) models.Task {
	var schedule models.Schedule
	err := db.C(models.Schedules).Find(bson.M{"ulid": ulid}).One(&schedule)
	if err == mgo.ErrNotFound {
		// TODO Register error
		return models.Task{}
	}

	var task models.Task
	if len(schedule.Tasks) == 1 {
		task = schedule.Tasks[0]
		// Remove the hole schedule due to no task remmaing
		err = db.C(models.Schedules).RemoveId(schedule.ID)
		if err != nil {
			// TODO Register error
			return models.Task{}
		}
	} else {
		for idx, t := range schedule.Tasks {
			if t.TaskID == taskID {
				task = t
				schedule.Tasks = append(schedule.Tasks[:idx], schedule.Tasks[idx+1:]...)
				err = db.C(models.Schedules).Update(bson.M{"ulid": ulid}, bson.M{"$set": bson.M{"tasks": schedule.Tasks}})
				if err != nil {
					// TODO Register error and handle the error
					log.Errorf("Error on extractTask for ULID: %s on TaskID: %s", ulid, taskID)
				} else {
					break
				}
			}
		}
	}
	return task
}

func (s *Server) updateTaskStatus(msg message.Message, status models.State, db *database.DataStore) error {
	selector := bson.M{"$and": []bson.M{
		bson.M{"ulid": msg.ULID},
		bson.M{"tasks": bson.M{"$elemMatch": bson.M{"task_id": msg.TaskID}}}}}
	update := bson.M{"$set": bson.M{"tasks.$.status": status}}
	if status == models.Failed {
		log.Errorf("[%s] %s", msg.ULID, errors.Errors[msg.ErrorID])
		s.saveErr(msg.ULID, msg.TaskID, msg.ErrorID, db)
	}
	return db.C(models.Schedules).Update(selector, update)
}

func newULID() ulid.ULID {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
