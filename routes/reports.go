package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	c "github.com/Azareal/Gosora/common"
	"github.com/Azareal/Gosora/common/counters"
)

func ReportSubmit(w http.ResponseWriter, r *http.Request, user c.User, sitemID string) c.RouteError {
	headerLite, ferr := c.SimpleUserCheck(w, r, &user)
	if ferr != nil {
		return ferr
	}
	isJs := (r.PostFormValue("isJs") == "1")

	itemID, err := strconv.Atoi(sitemID)
	if err != nil {
		return c.LocalError("Bad ID", w, r, user)
	}
	itemType := r.FormValue("type")

	// TODO: Localise these titles and bodies
	var title, content string
	if itemType == "reply" {
		reply, err := c.Rstore.Get(itemID)
		if err == sql.ErrNoRows {
			return c.LocalError("We were unable to find the reported post", w, r, user)
		} else if err != nil {
			return c.InternalError(err, w, r)
		}

		topic, err := c.Topics.Get(reply.ParentID)
		if err == sql.ErrNoRows {
			return c.LocalError("We weren't able to find the topic the reported post is supposed to be in", w, r, user)
		} else if err != nil {
			return c.InternalError(err, w, r)
		}

		title = "Reply: " + topic.Title
		content = reply.Content + "\n\nOriginal Post: #rid-" + strconv.Itoa(itemID)
	} else if itemType == "user-reply" {
		userReply, err := c.Prstore.Get(itemID)
		if err == sql.ErrNoRows {
			return c.LocalError("We weren't able to find the reported post", w, r, user)
		} else if err != nil {
			return c.InternalError(err, w, r)
		}

		profileOwner, err := c.Users.Get(userReply.ParentID)
		if err == sql.ErrNoRows {
			return c.LocalError("We weren't able to find the profile the reported post is supposed to be on", w, r, user)
		} else if err != nil {
			return c.InternalError(err, w, r)
		}
		title = "Profile: " + profileOwner.Name
		content = userReply.Content + "\n\nOriginal Post: @" + strconv.Itoa(userReply.ParentID)
	} else if itemType == "topic" {
		topic, err := c.Topics.Get(itemID)
		if err == sql.ErrNoRows {
			return c.NotFound(w, r, nil)
		} else if err != nil {
			return c.InternalError(err, w, r)
		}
		title = "Topic: " + topic.Title
		content = topic.Content + "\n\nOriginal Post: #tid-" + strconv.Itoa(itemID)
	} else {
		_, hasHook := headerLite.Hooks.VhookNeedHook("report_preassign", &itemID, &itemType)
		if hasHook {
			return nil
		}

		// Don't try to guess the type
		return c.LocalError("Unknown type", w, r, user)
	}

	// TODO: Repost attachments in the reports forum, so that the mods can see them
	_, err = c.Reports.Create(title, content, &user, itemType, itemID)
	if err == c.ErrAlreadyReported {
		return c.LocalError("Someone has already reported this!", w, r, user)
	}
	counters.PostCounter.Bump()

	if !isJs {
		// TODO: Redirect back to where we came from
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		_, _ = w.Write(successJSONBytes)
	}
	return nil
}
