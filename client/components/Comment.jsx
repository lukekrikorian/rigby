import React, { Component } from "react";
import { Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import "../css/Comment.css";
import "../css/Form.css";

export default class Comment extends Component {
	constructor(props){
		super(props);
		this.state = {
			replyBox: false,
			body: '',
			error: ''
		};
		this.change = this.change.bind(this);
		this.submit = this.submit.bind(this);
		this.toggle = this.toggle.bind(this);
	}

	toggle(event){
		event.preventDefault();
		this.setState({ replyBox: !this.state.replyBox });
	}

	change(event){
		const { target } = event;
		this.setState({
			[target.name]: target.value
		});
	}

	submit(event){
		event.preventDefault();

		const body = {
			body: this.state.body,
			parentID: this.props.Comment.ID
		};
		
		fetch("/api/replies", {
			method: "POST",
			credentials: "same-origin",
			body: JSON.stringify(body),
		}).then(response => {
			if(!response.ok) { 
				response.text().then(err => this.setState({ error: err })); 
			} else {
				location.reload();
			}
		}).catch(console.error);

	}

	render(){
		const comment = this.props.Comment;
		return (
			<div className={this.props.className || "comment"}>
				<p className="author"><Link to={`/~${comment.Author}`}>{comment.Author}</Link> on {this.props.showPost ? <Link to={`/posts/${comment.PostID}`}>this post</Link> : comment.Created.substring(0, 10)}</p>
				<ReactMarkdown allowedTypes={["paragraph", "text", "strong", "emphasis", "link"]} source={comment.Body} className="commentBody"/>
				<div className="replies">
				{ comment.Replies &&
					comment.Replies.map(reply => <Comment className="reply" Comment={reply}/>) }
				{ !this.props.className && <button class="linkLike" onClick={this.toggle}>Reply</button> }
				{ this.state.replyBox &&
					<div class="replyBox">
						<textarea style={{ margin: 8, width: "24em", marginLeft: 0 }} onChange={this.change} name="body" placeholder="Ok so basically..."></textarea>
						<br/>
						<button style={{ marginRight: 4 }} onClick={this.submit} class="linkLike">Send</button>
						<span>{this.state.error}</span>
					</div> }
				</div>
			</div>
		);
	}
}