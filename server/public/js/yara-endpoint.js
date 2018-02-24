
/*
    Functions for dashboard
*/

function renderDashboard() {
    $("#board-title").text("Dashboard");

    var tpl = _.unescape($("#dashboard-tpl").html());
    $("#page-body").empty();
    $("#page-body").html(tpl);
    getDashboard()

}

function updateOnline(online, total){
    Morris.Donut({
        data: [
            {label:"online", value: online},
            {label:"offline", value: total - online}
        ],
        element: "online-chart",
        resize: true,
    });
}

function updateMatches(matches, total){
    Morris.Donut({
        data: [
            {label:"clean", value: matches},
            {label:"infected", value: total - matches}
        ],
        color: "#FF0000",
        element: "infected-chart",
        resize: true,
    });
}

function updateEndpointsList(obj) {
    var tpl = ejs.compile(_.unescape($("#dashboard-list-of-assets-tpl").html()));
    var html = tpl({d: obj});
    $("#list-of-assets").empty();
    $("#list-of-assets").html(html);
    if ($("#list-of-assets").length >= 50) {
        $("#list-of-assets").find("tr td:nth-child(-n+50)").remove();
    }
}

function getDashboard(){
    $.getJSON("/dashboard", function(obj, status){
        if (status === "success") {
            if (_.isNull(obj.assets)) {
                $("#total-assets").text("0");
            } else {
                $("#total-assets").text(obj.assets.length);
            }
            if (_.isNull(obj.rules)) {
                $("#total-rules").text("0");
            } else {
                $("#total-rules").text(obj.rules.length);
            }

            var now = moment().subtract(5, 'minutes');
            var online = obj.assets.filter(function(o, i){
                return moment(o.last_ping) >  now
            });
            updateOnline(online.length, obj.assets.length);
            // updateMatches(online.length, obj.assets.length);
            updateEndpointsList(obj.assets);

        }
    });
}

/*
    Functions for assets
*/

function loadListOfAssets(event) {
    $("#board-title").text("List of assets");
    var tpl = _.unescape($("#list-of-assets-tpl").html());
    $("#page-body").empty();
    $("#page-body").html(tpl);
    getListOfAssets();
}

function updateListOfAssets(obj) {
    var tpl = ejs.compile(_.unescape($("#list-of-assets-list-tpl").html()));
    var tpl_modal = ejs.compile(_.unescape($("#list-of-assets-modal-tpl").html()));
    var html = tpl({d: obj});
    var html_modal = tpl_modal({d: obj});
    $("#list-of-assets").empty();
    $("#list-of-assets").html(html);
    $("#modals").empty();
    $("#modals").html(html_modal);
    hljs.initHighlighting();
}

function getListOfAssets() {
    $.getJSON("/assets", function(obj, status){
        if (status === "success") {
            updateListOfAssets(obj);
        }
    });
}


/*
    Functions for Rules
*/

function loadRules(event) {
    $("#board-title").text("Rules");
    var tpl = _.unescape($("#rules-tpl").html());
    $("#page-body").empty();
    $("#page-body").html(tpl);
    getRules();
}

function updateRules(obj) {
    var tpl = ejs.compile(_.unescape($("#rules-list-tpl").html()));
    var tpl_modal = ejs.compile(_.unescape($("#rules-modal-list-tpl").html()));
    var html = tpl({d: obj});
    var html_modal = tpl_modal({d: obj});
    $("#list-of-rules").empty();
    $("#list-of-rules").html(html);
    $("#modals").empty();
    $("#modals").html(html_modal);
    hljs.initHighlighting();
}

function loadNewRule(event) {
    $("#board-title").text("New Rule");
    var tpl = _.unescape($("#new-rule-tpl").html());
    $("#page-body").empty();
    $("#page-body").html(tpl);
    $("#submit-new-rule").submit(function(event){
        event.preventDefault();
        var data = {
            name: $("#new-rule-name").val(),
            tags: $("#new-rule-tags").val().split(","),
            data: $("#new-rule-data").val(),
        }
        $.ajax({
            type: 'POST',
            url:  '/rules/add',
            data:  JSON.stringify(data),
            dataType: 'json',
            contentType: 'application/json; charset=utf-8'
      }).always(function(data) {
        if (! _.isObject(data)) {
            data = JSON.parse(data);
        }
        if (_.has(data, "responseJSON") && data.responseJSON.error) {
            alert("The server was unable to insert the rule. Report this.\nErr: " + data.responseJSON.error_msg);
        } else if (data.error) {
            alert("The server was unable to insert the rule. Report this.\nErr: " + data.error_msg);
        } else {
            alert("Rule inserted correctly.")
            loadRules();
        }
      });
    });
}

function getRules() {
    $.getJSON("/rules", function(obj, status){
        if (status === "success") {
            updateRules(obj);
        }
    });
}

/*
    Functions for Tasks
*/

function loadTasks(event) {
    $("#board-title").text("Tasks");
    var tpl = _.unescape($("#tasks-tpl").html());
    $("#page-body").empty();
    $("#page-body").html(tpl);
    getTasks();
}

function updateTasks(obj) {
    var tpl = ejs.compile(_.unescape($("#tasks-list-tpl").html()));
    var tpl_modal = ejs.compile(_.unescape($("#task-modal-list-tpl").html()));
    var html = tpl({d: obj});
    var html_modal = tpl_modal({d: obj});
    $("#list-of-tasks").empty();
    $("#list-of-tasks").html(html);
    $("#modals").empty();
    $("#modals").html(html_modal);
}

function loadNewTask() {
    $("#board-title").text("New Schedule");
    var tpl = ejs.compile(_.unescape($("#new-task-tpl").html()));

    var assets = getAssetsList();
    var rules = getRulesList();
    var commands = getCommandsList();

    var data = {
        assets: assets,
        rules: rules,
        commands: commands,
    };

    var html = tpl({d: data});

    $("#page-body").empty();
    $("#page-body").html(html);

    $('.selectpicker').selectpicker('show');
    $('#new-task-datetime').datetimepicker({
        format: "DD/MM/YYYY hh:mm:ss a",
    });

    $("#submit-new-task").submit(function(event){
        event.preventDefault();
        var data = {
            assets: $("#new-task-asset").val(),
            rules: $("#new-task-rule").val(),
            command: $("#new-task-command").val(),
            target: $("#new-task-target").val(),
            when: $("#new-task-datetime").val(),
        }

        $.ajax({
        type: 'POST',
        url:  '/tasks/add',
        data:  JSON.stringify(data),
        dataType: 'json',
        contentType: 'application/json; charset=utf-8'
      }).always(function(data){
        if (! _.isObject(data)) {
            data = JSON.parse(data);
        }

        if (_.has(data, "responseJSON") && data.responseJSON.error) {
            alert("The server was unable to insert the task. Report this.\nErr: " + data.responseJSON.error_msg);
        } else if (data.error) {
            alert("The server was unable to insert the task. Report this.\nErr: " + data.error_msg);
        } else {
            alert("Task inserted correctly.")
            loadTasks();
        }
      });
    });
}

function getData(uri) {
    var result;
    $.ajax({
        async: false,
        url: uri,
        dataType: "json",
        success: function(data){
            result = data;
        }
    });
    return result;
}
function getAssetsList() {
    return getData("/assets")
}

function getRulesList() {
    return getData("/rules")
}

function getCommandsList() {
    return getData("/commands")
}

function getTasks() {
    $.getJSON("/tasks", function(obj, status){
        if (status === "success") {
            updateTasks(obj);
        }
    });
}


/*
    Functions for Reports
*/

function loadReports(event) {
    $("#board-title").text("Reports");
    var tpl = _.unescape($("#reports-tpl").html());
    $("#page-body").empty();
    $("#page-body").html(tpl);
    getReports();
}

function updateResults(obj) {
    // console.dir(obj);
    var tpl = ejs.compile(_.unescape($("#reports-list-tpl").html()));
    var tpl_modal = ejs.compile(_.unescape($("#reports-modal-list-tpl").html()));
    var html = tpl({d: obj});
    var html_modal = tpl_modal({d: obj});
    $("#list-of-reports").empty();
    $("#list-of-reports").html(html);
    $("#modals").empty();
    $("#modals").html(html_modal);
}

function getReports() {
    $.getJSON("/tasks/results", function(obj, status){
        if (status === "success") {
            updateResults(obj);
        }
    });
}

/*
    Auxiliar functions
*/

function remove_whitespaces(s){
    return s.replace( /\s/g, "")
}
