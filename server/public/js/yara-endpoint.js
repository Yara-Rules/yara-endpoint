
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
            $("#total-assets").text(obj.assets.length);
            $("#total-rules").text(obj.rules.length);

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
    var html = tpl({d: obj});
    $("#list-of-assets").empty();
    $("#list-of-assets").html(html);
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
    var html = tpl({d: obj});
    $("#list-of-rules").empty();
    $("#list-of-rules").html(html);
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
    var html = tpl({d: obj});
    $("#list-of-tasks").prepend(html);
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
    var html = tpl({d: obj});
    $("#list-of-reports").empty();
    $("#list-of-reports").html(html);
}

function getReports() {
    $.getJSON("/tasks/results", function(obj, status){
        if (status === "success") {
            updateResults(obj);
        }
    });
}
