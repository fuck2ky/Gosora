package main

import (
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"mime"
	"path/filepath"
	"io/ioutil"
	"html/template"
)

const hour int = 60 * 60
const day int = hour * 24
const month int = day * 30
const year int = day * 365
const kilobyte int = 1024
const megabyte int = 1024 * 1024
const saltLength int = 32
const sessionLength int = 80

var db *sql.DB
var get_session_stmt *sql.Stmt
var create_topic_stmt *sql.Stmt
var create_reply_stmt *sql.Stmt
var update_forum_cache_stmt *sql.Stmt
var edit_topic_stmt *sql.Stmt
var edit_reply_stmt *sql.Stmt
var delete_reply_stmt *sql.Stmt
var delete_topic_stmt *sql.Stmt
var stick_topic_stmt *sql.Stmt
var unstick_topic_stmt *sql.Stmt
var login_stmt *sql.Stmt
var update_session_stmt *sql.Stmt
var logout_stmt *sql.Stmt
var set_password_stmt *sql.Stmt
var get_password_stmt *sql.Stmt
var set_avatar_stmt *sql.Stmt
var set_username_stmt *sql.Stmt
var register_stmt *sql.Stmt
var username_exists_stmt *sql.Stmt
var create_profile_reply_stmt *sql.Stmt
var edit_profile_reply_stmt *sql.Stmt
var delete_profile_reply_stmt *sql.Stmt

var create_forum_stmt *sql.Stmt
var delete_forum_stmt *sql.Stmt
var update_forum_stmt *sql.Stmt

var custom_pages map[string]string = make(map[string]string)
var templates = template.Must(template.ParseGlob("templates/*"))
var no_css_tmpl = template.CSS("")
var staff_css_tmpl = template.CSS(staff_css)
var groups map[int]Group = make(map[int]Group)
var forums map[int]Forum = make(map[int]Forum)
var static_files map[string]SFile = make(map[string]SFile)

func init_database(err error) {
	if(dbpassword != ""){
		dbpassword = ":" + dbpassword
	}
	db, err = sql.Open("mysql",dbuser + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname)
	if err != nil {
		log.Fatal(err)
	}
	
	// Make sure that the connection is alive..
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing get_session statement.")
	get_session_stmt, err = db.Prepare("SELECT `uid`, `name`, `group`, `is_super_admin`, `session`, `avatar` FROM `users` WHERE `uid` = ? AND `session` = ? AND `session` <> ''")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing create_topic statement.")
	create_topic_stmt, err = db.Prepare("INSERT INTO topics(title,content,parsed_content,createdAt,createdBy) VALUES(?,?,?,NOW(),?)")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing create_reply statement.")
	create_reply_stmt, err = db.Prepare("INSERT INTO replies(tid,content,parsed_content,createdAt,createdBy) VALUES(?,?,?,NOW(),?)")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing update_forum_cache statement.")
	update_forum_cache_stmt, err = db.Prepare("UPDATE forums SET lastTopic = ?, lastTopicID = ?, lastReplyer = ?, lastReplyerID = ?, lastTopicTime = NOW() WHERE fid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing edit_topic statement.")
	edit_topic_stmt, err = db.Prepare("UPDATE topics SET title = ?, content = ?, parsed_content = ?, is_closed = ? WHERE tid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing edit_reply statement.")
	edit_reply_stmt, err = db.Prepare("UPDATE replies SET content = ?, parsed_content = ? WHERE rid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing delete_reply statement.")
	delete_reply_stmt, err = db.Prepare("DELETE FROM replies WHERE rid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing delete_topic statement.")
	delete_topic_stmt, err = db.Prepare("DELETE FROM topics WHERE tid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing stick_topic statement.")
	stick_topic_stmt, err = db.Prepare("UPDATE topics SET sticky = 1 WHERE tid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing unstick_topic statement.")
	unstick_topic_stmt, err = db.Prepare("UPDATE topics SET sticky = 0 WHERE tid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing login statement.")
	login_stmt, err = db.Prepare("SELECT `uid`, `name`, `password`, `salt` FROM `users` WHERE `name` = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing update_session statement.")
	update_session_stmt, err = db.Prepare("UPDATE users SET session = ? WHERE uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing logout statement.")
	logout_stmt, err = db.Prepare("UPDATE users SET session = '' WHERE uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing set_password statement.")
	set_password_stmt, err = db.Prepare("UPDATE users SET password = ?, salt = ? WHERE uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing get_password statement.")
	get_password_stmt, err = db.Prepare("SELECT `password`, `salt` FROM `users` WHERE `uid` = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing set_avatar statement.")
	set_avatar_stmt, err = db.Prepare("UPDATE users SET avatar = ? WHERE uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing set_username statement.")
	set_username_stmt, err = db.Prepare("UPDATE users SET name = ? WHERE uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	// Add an admin version of register_stmt with more flexibility
	// create_account_stmt, err = db.Prepare("INSERT INTO 
	
	log.Print("Preparing register statement.")
	register_stmt, err = db.Prepare("INSERT INTO users(`name`,`password`,`salt`,`group`,`is_super_admin`,`session`) VALUES(?,?,?,2,0,?)")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing username_exists statement.")
	username_exists_stmt, err = db.Prepare("SELECT `name` FROM `users` WHERE `name` = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing create_profile_reply statement.")
	create_profile_reply_stmt, err = db.Prepare("INSERT INTO users_replies(uid,content,parsed_content,createdAt,createdBy) VALUES(?,?,?,NOW(),?)")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing edit_profile_reply statement.")
	edit_profile_reply_stmt, err = db.Prepare("UPDATE users_replies SET content = ?, parsed_content = ? WHERE rid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing delete_profile_reply statement.")
	delete_profile_reply_stmt, err = db.Prepare("DELETE FROM users_replies WHERE rid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing create_forum statement.")
	create_forum_stmt, err = db.Prepare("INSERT INTO forums(name) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing delete_forum statement.")
	delete_forum_stmt, err = db.Prepare("DELETE FROM forums WHERE fid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Preparing update_forum statement.")
	update_forum_stmt, err = db.Prepare("UPDATE forums SET name = ? WHERE fid = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Loading the usergroups.")
	rows, err := db.Query("SELECT gid,name,permissions,is_admin,is_banned,tag FROM users_groups")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		group := Group{0,"","",false,false,""}
		err := rows.Scan(&group.ID, &group.Name, &group.Permissions, &group.Is_Admin, &group.Is_Banned, &group.Tag)
		if err != nil {
			log.Fatal(err)
		}
		groups[group.ID] = group
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Loading the forums.")
	rows, err = db.Query("SELECT fid, name, active, lastTopic, lastTopicID, lastReplyer, lastReplyerID, lastTopicTime FROM forums")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		forum := Forum{0,"",true,"",0,"",0,""}
		err := rows.Scan(&forum.ID, &forum.Name, &forum.Active, &forum.LastTopic, &forum.LastTopicID, &forum.LastReplyer, &forum.LastReplyerID, &forum.LastTopicTime)
		if err != nil {
			log.Fatal(err)
		}
		
		if forum.LastTopicID != 0 {
			forum.LastTopicTime, err = relative_time(forum.LastTopicTime)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			forum.LastTopic = "None"
			forum.LastTopicTime = ""
		}
		
		forums[forum.ID] = forum
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Adding the uncategorised forum")
	forums[0] = Forum{0,"Uncategorised",uncategorised_forum_visible,"",0,"",0,""}
	log.Print("Adding the reports forum")
	forums[-1] = Forum{-1,"Reports",false,"",0,"",0,""}
}

func main(){
	var err error
	init_database(err);
	
	log.Print("Loading the custom pages.")
	err = filepath.Walk("pages/", add_custom_page)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Print("Loading the static files.")
	files, err := ioutil.ReadDir("./public")
	if err != nil {
		log.Fatal(err)
	}
	
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		data, err := ioutil.ReadFile("./public/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		
		log.Print("Added the '" + f.Name() + "' static file.")
		static_files["/static/" + f.Name()] = SFile{data,0,int64(len(data)),mime.TypeByExtension(filepath.Ext(f.Name())),f,f.ModTime().UTC().Format(http.TimeFormat)}
	}
	
	// In a directory to stop it clashing with the other paths
	http.HandleFunc("/static/", route_static)
	//http.HandleFunc("/static/", route_fstatic)
	//fs_p := http.FileServer(http.Dir("./public"))
	//http.Handle("/static/", http.StripPrefix("/static/",fs_p))
	fs_u := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/",fs_u))
	
	http.HandleFunc("/overview/", route_overview)
	http.HandleFunc("/topics/create/", route_topic_create)
	http.HandleFunc("/topics/", route_topics)
	http.HandleFunc("/forums/", route_forums)
	http.HandleFunc("/forum/", route_forum)
	http.HandleFunc("/topic/create/submit/", route_create_topic) //POST
	http.HandleFunc("/topic/", route_topic_id)
	http.HandleFunc("/reply/create/", route_create_reply) //POST
	//http.HandleFunc("/reply/edit/", route_reply_edit) //POST
	//http.HandleFunc("/reply/delete/", route_reply_delete) //POST
	http.HandleFunc("/reply/edit/submit/", route_reply_edit_submit) //POST
	http.HandleFunc("/reply/delete/submit/", route_reply_delete_submit) //POST
	http.HandleFunc("/topic/edit/submit/", route_edit_topic) //POST
	http.HandleFunc("/topic/delete/submit/", route_delete_topic)
	http.HandleFunc("/topic/stick/submit/", route_stick_topic)
	http.HandleFunc("/topic/unstick/submit/", route_unstick_topic)
	
	// Custom Pages
	http.HandleFunc("/pages/", route_custom_page)
	
	// Accounts
	http.HandleFunc("/accounts/login/", route_login)
	http.HandleFunc("/accounts/create/", route_register)
	http.HandleFunc("/accounts/logout/", route_logout)
	http.HandleFunc("/accounts/login/submit/", route_login_submit) // POST
	http.HandleFunc("/accounts/create/submit/", route_register_submit) // POST
	
	//http.HandleFunc("/accounts/list/", route_login) // Redirect /accounts/ and /user/ to here..
	//http.HandleFunc("/accounts/create/full/", route_logout)
	//http.HandleFunc("/user/edit/", route_logout)
	http.HandleFunc("/user/edit/critical/", route_account_own_edit_critical) // Password & Email
	http.HandleFunc("/user/edit/critical/submit/", route_account_own_edit_critical_submit)
	http.HandleFunc("/user/edit/avatar/", route_account_own_edit_avatar)
	http.HandleFunc("/user/edit/avatar/submit/", route_account_own_edit_avatar_submit)
	http.HandleFunc("/user/edit/username/", route_account_own_edit_username)
	http.HandleFunc("/user/edit/username/submit/", route_account_own_edit_username_submit)
	http.HandleFunc("/user/", route_profile)
	http.HandleFunc("/profile/reply/create/", route_profile_reply_create)
	http.HandleFunc("/profile/reply/edit/submit/", route_profile_reply_edit_submit)
	http.HandleFunc("/profile/reply/delete/submit/", route_profile_reply_delete_submit)
	//http.HandleFunc("/user/:id/edit/", route_logout)
	//http.HandleFunc("/user/:id/ban/", route_logout)
	
	// Admin
	http.HandleFunc("/panel/forums/", route_panel_forums)
	http.HandleFunc("/panel/forums/create/", route_panel_forums_create_submit)
	http.HandleFunc("/panel/forums/delete/", route_panel_forums_delete)
	http.HandleFunc("/panel/forums/delete/submit/", route_panel_forums_delete_submit)
	http.HandleFunc("/panel/forums/edit/submit/", route_panel_forums_edit_submit)
	
	http.HandleFunc("/", default_route)
	
	defer db.Close()
    http.ListenAndServe(":8080", nil)
}