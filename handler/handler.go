package handler

import (
	"encoding/json"

	"github.com/bentranter/lazyblog/model"
	"github.com/dinever/golf"
)

var s = model.NewMemoryStore()

// Info displays info about the available API routes.
//
// @TODO: Don't marshal JSON for every request, just hard code this as a byte
//        slice.
func Info(ctx *golf.Context) {
	routes := map[string]string{
		"posts_url": "/api/posts{/id}",
	}
	resp, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		ctx.Abort(500)
		return
	}
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.Write(resp)
}

// SetPost is the API method for creating a new post.
func SetPost(ctx *golf.Context) {
	v := &struct {
		Title string
		Body  string
	}{}
	err := json.NewDecoder(ctx.Request.Body).Decode(v)
	if err != nil {
		ctx.Abort(500)
		return
	}
	defer ctx.Request.Body.Close()

	p := model.NewPost(v.Title, v.Body)
	id, err := p.Set(s)
	if err != nil {
		ctx.Abort(500)
		return
	}
	ctx.Response.Header().Set("Content-Type", "application/json")
	resp, err := json.MarshalIndent(map[string]string{
		"id": string(id),
	}, "", "  ")
	if err != nil {
		ctx.Abort(500)
	}
	ctx.Response.Write(resp)
}

// GetPostJSON is the API method for getting a post's JSON.
func GetPostJSON(ctx *golf.Context) {
	id := ctx.Param("id")
	postJSON, err := model.GetJSON([]byte(id), s)
	if err != nil {
		ctx.Abort(404)
		return
	}
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.Write(postJSON)
}

// GetAllPosts is a method for getting every post
func GetAllPosts(ctx *golf.Context) {
	posts, err := model.GetAll(s)
	if err != nil {
		ctx.Abort(404)
		return
	}

	resp, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		ctx.Abort(500)
		return
	}
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.Write(resp)
}
