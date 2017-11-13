// +build mssql

// This file was generated by Gosora's Query Generator. Please try to avoid modifying this file, as it might change at any time.
package main

import "log"
import "database/sql"
import "./common"

// nolint
type Stmts struct {
	getPassword *sql.Stmt
	getSettings *sql.Stmt
	getSetting *sql.Stmt
	getFullSetting *sql.Stmt
	isPluginActive *sql.Stmt
	getUsersOffset *sql.Stmt
	isThemeDefault *sql.Stmt
	getModlogs *sql.Stmt
	getModlogsOffset *sql.Stmt
	getReplyTID *sql.Stmt
	getTopicFID *sql.Stmt
	getUserReplyUID *sql.Stmt
	getUserName *sql.Stmt
	getEmailsByUser *sql.Stmt
	getTopicBasic *sql.Stmt
	getActivityEntry *sql.Stmt
	forumEntryExists *sql.Stmt
	groupEntryExists *sql.Stmt
	getForumTopicsOffset *sql.Stmt
	getAttachment *sql.Stmt
	getTopicRepliesOffset *sql.Stmt
	getTopicList *sql.Stmt
	getTopicReplies *sql.Stmt
	getForumTopics *sql.Stmt
	getProfileReplies *sql.Stmt
	getWatchers *sql.Stmt
	createReport *sql.Stmt
	addActivity *sql.Stmt
	notifyOne *sql.Stmt
	addEmail *sql.Stmt
	addSubscription *sql.Stmt
	addForumPermsToForum *sql.Stmt
	addPlugin *sql.Stmt
	addTheme *sql.Stmt
	addAttachment *sql.Stmt
	createWordFilter *sql.Stmt
	editReply *sql.Stmt
	editProfileReply *sql.Stmt
	updateSetting *sql.Stmt
	updatePlugin *sql.Stmt
	updatePluginInstall *sql.Stmt
	updateTheme *sql.Stmt
	updateUser *sql.Stmt
	updateGroupPerms *sql.Stmt
	updateGroup *sql.Stmt
	updateEmail *sql.Stmt
	verifyEmail *sql.Stmt
	setTempGroup *sql.Stmt
	updateWordFilter *sql.Stmt
	bumpSync *sql.Stmt
	deleteProfileReply *sql.Stmt
	deleteActivityStreamMatch *sql.Stmt
	deleteWordFilter *sql.Stmt
	reportExists *sql.Stmt
	modlogCount *sql.Stmt
	notifyWatchers *sql.Stmt

	getActivityFeedByWatcher *sql.Stmt
	getActivityCountByWatcher *sql.Stmt
	todaysPostCount *sql.Stmt
	todaysTopicCount *sql.Stmt
	todaysReportCount *sql.Stmt
	todaysNewUserCount *sql.Stmt
	findUsersByIPUsers *sql.Stmt
	findUsersByIPTopics *sql.Stmt
	findUsersByIPReplies *sql.Stmt

	Mocks bool
}

// nolint
func _gen_mssql() (err error) {
	if common.Dev.DebugMode {
		log.Print("Building the generated statements")
	}
	
	log.Print("Preparing getPassword statement.")
	stmts.getPassword, err = db.Prepare("SELECT [password],[salt] FROM [users] WHERE [uid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [password],[salt] FROM [users] WHERE [uid] = ?1")
		return err
	}
		
	log.Print("Preparing getSettings statement.")
	stmts.getSettings, err = db.Prepare("SELECT [name],[content],[type] FROM [settings]")
	if err != nil {
		log.Print("Bad Query: ","SELECT [name],[content],[type] FROM [settings]")
		return err
	}
		
	log.Print("Preparing getSetting statement.")
	stmts.getSetting, err = db.Prepare("SELECT [content],[type] FROM [settings] WHERE [name] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [content],[type] FROM [settings] WHERE [name] = ?1")
		return err
	}
		
	log.Print("Preparing getFullSetting statement.")
	stmts.getFullSetting, err = db.Prepare("SELECT [name],[type],[constraints] FROM [settings] WHERE [name] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [name],[type],[constraints] FROM [settings] WHERE [name] = ?1")
		return err
	}
		
	log.Print("Preparing isPluginActive statement.")
	stmts.isPluginActive, err = db.Prepare("SELECT [active] FROM [plugins] WHERE [uname] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [active] FROM [plugins] WHERE [uname] = ?1")
		return err
	}
		
	log.Print("Preparing getUsersOffset statement.")
	stmts.getUsersOffset, err = db.Prepare("SELECT [uid],[name],[group],[active],[is_super_admin],[avatar] FROM [users] ORDER BY uid ASC OFFSET ?1 ROWS FETCH NEXT ?2 ROWS ONLY")
	if err != nil {
		log.Print("Bad Query: ","SELECT [uid],[name],[group],[active],[is_super_admin],[avatar] FROM [users] ORDER BY uid ASC OFFSET ?1 ROWS FETCH NEXT ?2 ROWS ONLY")
		return err
	}
		
	log.Print("Preparing isThemeDefault statement.")
	stmts.isThemeDefault, err = db.Prepare("SELECT [default] FROM [themes] WHERE [uname] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [default] FROM [themes] WHERE [uname] = ?1")
		return err
	}
		
	log.Print("Preparing getModlogs statement.")
	stmts.getModlogs, err = db.Prepare("SELECT [action],[elementID],[elementType],[ipaddress],[actorID],[doneAt] FROM [moderation_logs]")
	if err != nil {
		log.Print("Bad Query: ","SELECT [action],[elementID],[elementType],[ipaddress],[actorID],[doneAt] FROM [moderation_logs]")
		return err
	}
		
	log.Print("Preparing getModlogsOffset statement.")
	stmts.getModlogsOffset, err = db.Prepare("SELECT [action],[elementID],[elementType],[ipaddress],[actorID],[doneAt] FROM [moderation_logs] ORDER BY doneAt DESC OFFSET ?1 ROWS FETCH NEXT ?2 ROWS ONLY")
	if err != nil {
		log.Print("Bad Query: ","SELECT [action],[elementID],[elementType],[ipaddress],[actorID],[doneAt] FROM [moderation_logs] ORDER BY doneAt DESC OFFSET ?1 ROWS FETCH NEXT ?2 ROWS ONLY")
		return err
	}
		
	log.Print("Preparing getReplyTID statement.")
	stmts.getReplyTID, err = db.Prepare("SELECT [tid] FROM [replies] WHERE [rid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [tid] FROM [replies] WHERE [rid] = ?1")
		return err
	}
		
	log.Print("Preparing getTopicFID statement.")
	stmts.getTopicFID, err = db.Prepare("SELECT [parentID] FROM [topics] WHERE [tid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [parentID] FROM [topics] WHERE [tid] = ?1")
		return err
	}
		
	log.Print("Preparing getUserReplyUID statement.")
	stmts.getUserReplyUID, err = db.Prepare("SELECT [uid] FROM [users_replies] WHERE [rid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [uid] FROM [users_replies] WHERE [rid] = ?1")
		return err
	}
		
	log.Print("Preparing getUserName statement.")
	stmts.getUserName, err = db.Prepare("SELECT [name] FROM [users] WHERE [uid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [name] FROM [users] WHERE [uid] = ?1")
		return err
	}
		
	log.Print("Preparing getEmailsByUser statement.")
	stmts.getEmailsByUser, err = db.Prepare("SELECT [email],[validated],[token] FROM [emails] WHERE [uid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [email],[validated],[token] FROM [emails] WHERE [uid] = ?1")
		return err
	}
		
	log.Print("Preparing getTopicBasic statement.")
	stmts.getTopicBasic, err = db.Prepare("SELECT [title],[content] FROM [topics] WHERE [tid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [title],[content] FROM [topics] WHERE [tid] = ?1")
		return err
	}
		
	log.Print("Preparing getActivityEntry statement.")
	stmts.getActivityEntry, err = db.Prepare("SELECT [actor],[targetUser],[event],[elementType],[elementID] FROM [activity_stream] WHERE [asid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [actor],[targetUser],[event],[elementType],[elementID] FROM [activity_stream] WHERE [asid] = ?1")
		return err
	}
		
	log.Print("Preparing forumEntryExists statement.")
	stmts.forumEntryExists, err = db.Prepare("SELECT [fid] FROM [forums] WHERE [name] = '' ORDER BY fid ASC OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY")
	if err != nil {
		log.Print("Bad Query: ","SELECT [fid] FROM [forums] WHERE [name] = '' ORDER BY fid ASC OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY")
		return err
	}
		
	log.Print("Preparing groupEntryExists statement.")
	stmts.groupEntryExists, err = db.Prepare("SELECT [gid] FROM [users_groups] WHERE [name] = '' ORDER BY gid ASC OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY")
	if err != nil {
		log.Print("Bad Query: ","SELECT [gid] FROM [users_groups] WHERE [name] = '' ORDER BY gid ASC OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY")
		return err
	}
		
	log.Print("Preparing getForumTopicsOffset statement.")
	stmts.getForumTopicsOffset, err = db.Prepare("SELECT [tid],[title],[content],[createdBy],[is_closed],[sticky],[createdAt],[lastReplyAt],[lastReplyBy],[parentID],[postCount],[likeCount] FROM [topics] WHERE [parentID] = ?1 ORDER BY sticky DESC,lastReplyAt DESC,createdBy DESC OFFSET ?2 ROWS FETCH NEXT ?3 ROWS ONLY")
	if err != nil {
		log.Print("Bad Query: ","SELECT [tid],[title],[content],[createdBy],[is_closed],[sticky],[createdAt],[lastReplyAt],[lastReplyBy],[parentID],[postCount],[likeCount] FROM [topics] WHERE [parentID] = ?1 ORDER BY sticky DESC,lastReplyAt DESC,createdBy DESC OFFSET ?2 ROWS FETCH NEXT ?3 ROWS ONLY")
		return err
	}
		
	log.Print("Preparing getAttachment statement.")
	stmts.getAttachment, err = db.Prepare("SELECT [sectionID],[sectionTable],[originID],[originTable],[uploadedBy],[path] FROM [attachments] WHERE [path] = ?1 AND [sectionID] = ?2 AND [sectionTable] = ?3")
	if err != nil {
		log.Print("Bad Query: ","SELECT [sectionID],[sectionTable],[originID],[originTable],[uploadedBy],[path] FROM [attachments] WHERE [path] = ?1 AND [sectionID] = ?2 AND [sectionTable] = ?3")
		return err
	}
		
	log.Print("Preparing getTopicRepliesOffset statement.")
	stmts.getTopicRepliesOffset, err = db.Prepare("SELECT [replies].[rid],[replies].[content],[replies].[createdBy],[replies].[createdAt],[replies].[lastEdit],[replies].[lastEditBy],[users].[avatar],[users].[name],[users].[group],[users].[url_prefix],[users].[url_name],[users].[level],[replies].[ipaddress],[replies].[likeCount],[replies].[actionType] FROM [replies] LEFT JOIN [users] ON [replies].[createdBy] = [users].[uid]  WHERE [replies].[tid] = ?1 ORDER BY replies.rid ASC OFFSET ?2 ROWS FETCH NEXT ?3 ROWS ONLY")
	if err != nil {
		log.Print("Bad Query: ","SELECT [replies].[rid],[replies].[content],[replies].[createdBy],[replies].[createdAt],[replies].[lastEdit],[replies].[lastEditBy],[users].[avatar],[users].[name],[users].[group],[users].[url_prefix],[users].[url_name],[users].[level],[replies].[ipaddress],[replies].[likeCount],[replies].[actionType] FROM [replies] LEFT JOIN [users] ON [replies].[createdBy] = [users].[uid]  WHERE [replies].[tid] = ?1 ORDER BY replies.rid ASC OFFSET ?2 ROWS FETCH NEXT ?3 ROWS ONLY")
		return err
	}
		
	log.Print("Preparing getTopicList statement.")
	stmts.getTopicList, err = db.Prepare("SELECT [topics].[tid],[topics].[title],[topics].[content],[topics].[createdBy],[topics].[is_closed],[topics].[sticky],[topics].[createdAt],[topics].[parentID],[users].[name],[users].[avatar] FROM [topics] LEFT JOIN [users] ON [topics].[createdBy] = [users].[uid]  ORDER BY topics.sticky DESC,topics.lastReplyAt DESC,topics.createdBy DESC")
	if err != nil {
		log.Print("Bad Query: ","SELECT [topics].[tid],[topics].[title],[topics].[content],[topics].[createdBy],[topics].[is_closed],[topics].[sticky],[topics].[createdAt],[topics].[parentID],[users].[name],[users].[avatar] FROM [topics] LEFT JOIN [users] ON [topics].[createdBy] = [users].[uid]  ORDER BY topics.sticky DESC,topics.lastReplyAt DESC,topics.createdBy DESC")
		return err
	}
		
	log.Print("Preparing getTopicReplies statement.")
	stmts.getTopicReplies, err = db.Prepare("SELECT [replies].[rid],[replies].[content],[replies].[createdBy],[replies].[createdAt],[replies].[lastEdit],[replies].[lastEditBy],[users].[avatar],[users].[name],[users].[group],[users].[url_prefix],[users].[url_name],[users].[level],[replies].[ipaddress] FROM [replies] LEFT JOIN [users] ON [replies].[createdBy] = [users].[uid]  WHERE [tid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [replies].[rid],[replies].[content],[replies].[createdBy],[replies].[createdAt],[replies].[lastEdit],[replies].[lastEditBy],[users].[avatar],[users].[name],[users].[group],[users].[url_prefix],[users].[url_name],[users].[level],[replies].[ipaddress] FROM [replies] LEFT JOIN [users] ON [replies].[createdBy] = [users].[uid]  WHERE [tid] = ?1")
		return err
	}
		
	log.Print("Preparing getForumTopics statement.")
	stmts.getForumTopics, err = db.Prepare("SELECT [topics].[tid],[topics].[title],[topics].[content],[topics].[createdBy],[topics].[is_closed],[topics].[sticky],[topics].[createdAt],[topics].[lastReplyAt],[topics].[parentID],[users].[name],[users].[avatar] FROM [topics] LEFT JOIN [users] ON [topics].[createdBy] = [users].[uid]  WHERE [topics].[parentID] = ?1 ORDER BY topics.sticky DESC,topics.lastReplyAt DESC,topics.createdBy DESC")
	if err != nil {
		log.Print("Bad Query: ","SELECT [topics].[tid],[topics].[title],[topics].[content],[topics].[createdBy],[topics].[is_closed],[topics].[sticky],[topics].[createdAt],[topics].[lastReplyAt],[topics].[parentID],[users].[name],[users].[avatar] FROM [topics] LEFT JOIN [users] ON [topics].[createdBy] = [users].[uid]  WHERE [topics].[parentID] = ?1 ORDER BY topics.sticky DESC,topics.lastReplyAt DESC,topics.createdBy DESC")
		return err
	}
		
	log.Print("Preparing getProfileReplies statement.")
	stmts.getProfileReplies, err = db.Prepare("SELECT [users_replies].[rid],[users_replies].[content],[users_replies].[createdBy],[users_replies].[createdAt],[users_replies].[lastEdit],[users_replies].[lastEditBy],[users].[avatar],[users].[name],[users].[group] FROM [users_replies] LEFT JOIN [users] ON [users_replies].[createdBy] = [users].[uid]  WHERE [users_replies].[uid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [users_replies].[rid],[users_replies].[content],[users_replies].[createdBy],[users_replies].[createdAt],[users_replies].[lastEdit],[users_replies].[lastEditBy],[users].[avatar],[users].[name],[users].[group] FROM [users_replies] LEFT JOIN [users] ON [users_replies].[createdBy] = [users].[uid]  WHERE [users_replies].[uid] = ?1")
		return err
	}
		
	log.Print("Preparing getWatchers statement.")
	stmts.getWatchers, err = db.Prepare("SELECT [activity_subscriptions].[user] FROM [activity_stream] INNER JOIN [activity_subscriptions] ON [activity_subscriptions].[targetType] = [activity_stream].[elementType] AND [activity_subscriptions].[targetID] = [activity_stream].[elementID] AND [activity_subscriptions].[user] != [activity_stream].[actor]  WHERE [asid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","SELECT [activity_subscriptions].[user] FROM [activity_stream] INNER JOIN [activity_subscriptions] ON [activity_subscriptions].[targetType] = [activity_stream].[elementType] AND [activity_subscriptions].[targetID] = [activity_stream].[elementID] AND [activity_subscriptions].[user] != [activity_stream].[actor]  WHERE [asid] = ?1")
		return err
	}
		
	log.Print("Preparing createReport statement.")
	stmts.createReport, err = db.Prepare("INSERT INTO [topics] ([title],[content],[parsed_content],[createdAt],[lastReplyAt],[createdBy],[lastReplyBy],[data],[parentID],[css_class]) VALUES (?,?,?,GETUTCDATE(),GETUTCDATE(),?,?,?,1,'report')")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [topics] ([title],[content],[parsed_content],[createdAt],[lastReplyAt],[createdBy],[lastReplyBy],[data],[parentID],[css_class]) VALUES (?,?,?,GETUTCDATE(),GETUTCDATE(),?,?,?,1,'report')")
		return err
	}
		
	log.Print("Preparing addActivity statement.")
	stmts.addActivity, err = db.Prepare("INSERT INTO [activity_stream] ([actor],[targetUser],[event],[elementType],[elementID]) VALUES (?,?,?,?,?)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [activity_stream] ([actor],[targetUser],[event],[elementType],[elementID]) VALUES (?,?,?,?,?)")
		return err
	}
		
	log.Print("Preparing notifyOne statement.")
	stmts.notifyOne, err = db.Prepare("INSERT INTO [activity_stream_matches] ([watcher],[asid]) VALUES (?,?)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [activity_stream_matches] ([watcher],[asid]) VALUES (?,?)")
		return err
	}
		
	log.Print("Preparing addEmail statement.")
	stmts.addEmail, err = db.Prepare("INSERT INTO [emails] ([email],[uid],[validated],[token]) VALUES (?,?,?,?)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [emails] ([email],[uid],[validated],[token]) VALUES (?,?,?,?)")
		return err
	}
		
	log.Print("Preparing addSubscription statement.")
	stmts.addSubscription, err = db.Prepare("INSERT INTO [activity_subscriptions] ([user],[targetID],[targetType],[level]) VALUES (?,?,?,2)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [activity_subscriptions] ([user],[targetID],[targetType],[level]) VALUES (?,?,?,2)")
		return err
	}
		
	log.Print("Preparing addForumPermsToForum statement.")
	stmts.addForumPermsToForum, err = db.Prepare("INSERT INTO [forums_permissions] ([gid],[fid],[preset],[permissions]) VALUES (?,?,?,?)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [forums_permissions] ([gid],[fid],[preset],[permissions]) VALUES (?,?,?,?)")
		return err
	}
		
	log.Print("Preparing addPlugin statement.")
	stmts.addPlugin, err = db.Prepare("INSERT INTO [plugins] ([uname],[active],[installed]) VALUES (?,?,?)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [plugins] ([uname],[active],[installed]) VALUES (?,?,?)")
		return err
	}
		
	log.Print("Preparing addTheme statement.")
	stmts.addTheme, err = db.Prepare("INSERT INTO [themes] ([uname],[default]) VALUES (?,?)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [themes] ([uname],[default]) VALUES (?,?)")
		return err
	}
		
	log.Print("Preparing addAttachment statement.")
	stmts.addAttachment, err = db.Prepare("INSERT INTO [attachments] ([sectionID],[sectionTable],[originID],[originTable],[uploadedBy],[path]) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [attachments] ([sectionID],[sectionTable],[originID],[originTable],[uploadedBy],[path]) VALUES (?,?,?,?,?,?)")
		return err
	}
		
	log.Print("Preparing createWordFilter statement.")
	stmts.createWordFilter, err = db.Prepare("INSERT INTO [word_filters] ([find],[replacement]) VALUES (?,?)")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [word_filters] ([find],[replacement]) VALUES (?,?)")
		return err
	}
		
	log.Print("Preparing editReply statement.")
	stmts.editReply, err = db.Prepare("UPDATE [replies] SET [content] = ?,[parsed_content] = ? WHERE [rid] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [replies] SET [content] = ?,[parsed_content] = ? WHERE [rid] = ?")
		return err
	}
		
	log.Print("Preparing editProfileReply statement.")
	stmts.editProfileReply, err = db.Prepare("UPDATE [users_replies] SET [content] = ?,[parsed_content] = ? WHERE [rid] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [users_replies] SET [content] = ?,[parsed_content] = ? WHERE [rid] = ?")
		return err
	}
		
	log.Print("Preparing updateSetting statement.")
	stmts.updateSetting, err = db.Prepare("UPDATE [settings] SET [content] = ? WHERE [name] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [settings] SET [content] = ? WHERE [name] = ?")
		return err
	}
		
	log.Print("Preparing updatePlugin statement.")
	stmts.updatePlugin, err = db.Prepare("UPDATE [plugins] SET [active] = ? WHERE [uname] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [plugins] SET [active] = ? WHERE [uname] = ?")
		return err
	}
		
	log.Print("Preparing updatePluginInstall statement.")
	stmts.updatePluginInstall, err = db.Prepare("UPDATE [plugins] SET [installed] = ? WHERE [uname] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [plugins] SET [installed] = ? WHERE [uname] = ?")
		return err
	}
		
	log.Print("Preparing updateTheme statement.")
	stmts.updateTheme, err = db.Prepare("UPDATE [themes] SET [default] = ? WHERE [uname] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [themes] SET [default] = ? WHERE [uname] = ?")
		return err
	}
		
	log.Print("Preparing updateUser statement.")
	stmts.updateUser, err = db.Prepare("UPDATE [users] SET [name] = ?,[email] = ?,[group] = ? WHERE [uid] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [users] SET [name] = ?,[email] = ?,[group] = ? WHERE [uid] = ?")
		return err
	}
		
	log.Print("Preparing updateGroupPerms statement.")
	stmts.updateGroupPerms, err = db.Prepare("UPDATE [users_groups] SET [permissions] = ? WHERE [gid] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [users_groups] SET [permissions] = ? WHERE [gid] = ?")
		return err
	}
		
	log.Print("Preparing updateGroup statement.")
	stmts.updateGroup, err = db.Prepare("UPDATE [users_groups] SET [name] = ?,[tag] = ? WHERE [gid] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [users_groups] SET [name] = ?,[tag] = ? WHERE [gid] = ?")
		return err
	}
		
	log.Print("Preparing updateEmail statement.")
	stmts.updateEmail, err = db.Prepare("UPDATE [emails] SET [email] = ?,[uid] = ?,[validated] = ?,[token] = ? WHERE [email] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [emails] SET [email] = ?,[uid] = ?,[validated] = ?,[token] = ? WHERE [email] = ?")
		return err
	}
		
	log.Print("Preparing verifyEmail statement.")
	stmts.verifyEmail, err = db.Prepare("UPDATE [emails] SET [validated] = 1,[token] = '' WHERE [email] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [emails] SET [validated] = 1,[token] = '' WHERE [email] = ?")
		return err
	}
		
	log.Print("Preparing setTempGroup statement.")
	stmts.setTempGroup, err = db.Prepare("UPDATE [users] SET [temp_group] = ? WHERE [uid] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [users] SET [temp_group] = ? WHERE [uid] = ?")
		return err
	}
		
	log.Print("Preparing updateWordFilter statement.")
	stmts.updateWordFilter, err = db.Prepare("UPDATE [word_filters] SET [find] = ?,[replacement] = ? WHERE [wfid] = ?")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [word_filters] SET [find] = ?,[replacement] = ? WHERE [wfid] = ?")
		return err
	}
		
	log.Print("Preparing bumpSync statement.")
	stmts.bumpSync, err = db.Prepare("UPDATE [sync] SET [last_update] = GETUTCDATE()")
	if err != nil {
		log.Print("Bad Query: ","UPDATE [sync] SET [last_update] = GETUTCDATE()")
		return err
	}
		
	log.Print("Preparing deleteProfileReply statement.")
	stmts.deleteProfileReply, err = db.Prepare("DELETE FROM [users_replies] WHERE [rid] = ?")
	if err != nil {
		log.Print("Bad Query: ","DELETE FROM [users_replies] WHERE [rid] = ?")
		return err
	}
		
	log.Print("Preparing deleteActivityStreamMatch statement.")
	stmts.deleteActivityStreamMatch, err = db.Prepare("DELETE FROM [activity_stream_matches] WHERE [watcher] = ? AND [asid] = ?")
	if err != nil {
		log.Print("Bad Query: ","DELETE FROM [activity_stream_matches] WHERE [watcher] = ? AND [asid] = ?")
		return err
	}
		
	log.Print("Preparing deleteWordFilter statement.")
	stmts.deleteWordFilter, err = db.Prepare("DELETE FROM [word_filters] WHERE [wfid] = ?")
	if err != nil {
		log.Print("Bad Query: ","DELETE FROM [word_filters] WHERE [wfid] = ?")
		return err
	}
		
	log.Print("Preparing reportExists statement.")
	stmts.reportExists, err = db.Prepare("SELECT COUNT(*) AS [count] FROM [topics] WHERE [data] = ? AND [data] != '' AND [parentID] = 1")
	if err != nil {
		log.Print("Bad Query: ","SELECT COUNT(*) AS [count] FROM [topics] WHERE [data] = ? AND [data] != '' AND [parentID] = 1")
		return err
	}
		
	log.Print("Preparing modlogCount statement.")
	stmts.modlogCount, err = db.Prepare("SELECT COUNT(*) AS [count] FROM [moderation_logs]")
	if err != nil {
		log.Print("Bad Query: ","SELECT COUNT(*) AS [count] FROM [moderation_logs]")
		return err
	}
		
	log.Print("Preparing notifyWatchers statement.")
	stmts.notifyWatchers, err = db.Prepare("INSERT INTO [activity_stream_matches] ([watcher],[asid]) SELECT [activity_subscriptions].[user],[activity_stream].[asid] FROM [activity_stream] INNER JOIN [activity_subscriptions] ON [activity_subscriptions].[targetType] = [activity_stream].[elementType] AND [activity_subscriptions].[targetID] = [activity_stream].[elementID] AND [activity_subscriptions].[user] != [activity_stream].[actor]  WHERE [asid] = ?1")
	if err != nil {
		log.Print("Bad Query: ","INSERT INTO [activity_stream_matches] ([watcher],[asid]) SELECT [activity_subscriptions].[user],[activity_stream].[asid] FROM [activity_stream] INNER JOIN [activity_subscriptions] ON [activity_subscriptions].[targetType] = [activity_stream].[elementType] AND [activity_subscriptions].[targetID] = [activity_stream].[elementID] AND [activity_subscriptions].[user] != [activity_stream].[actor]  WHERE [asid] = ?1")
		return err
	}
	
	return nil
}
