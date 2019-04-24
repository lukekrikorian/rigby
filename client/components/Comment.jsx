import React, { Component } from "react";
import { Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";

export default class Comment extends Component {
	constructor(props){
		super(props);
		this.state = {
			body: '',
			error: '',
			showCommentBox: false,
		};
		this.change = this.change.bind(this);
		this.submit = this.submit.bind(this);
	}

	change(event){
		const { target } = event;
		this.setState({
			[target.name]: target.value
		});
	}

	submit(event){
		event.preventDefault();

		if (!this.state.showCommentBox) {
			this.setState({ showCommentBox: true });
			return;
		}

		const body = {
			body: this.state.body,
			parentID: this.props.Comment.ID,
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
		const { Author, Body, Replies } = this.props.Comment;
		return (
			<div className={this.props.reply ? "reply" : "comment"}>
				<Link 
					to={`/~${Author}`}
					className="commentAuthor">
						~{ Author }
				</Link>
				<ReactMarkdown
					allowedTypes={["paragraph", "text", "strong", "emphasis", "link"]} 
					className="commentBody">
						{ Body }
				</ReactMarkdown>
				{ !this.props.reply &&
					<div className="repliesWrap">
						<div className="replies">
						{ Replies && Replies.map(reply => <Comment Comment={reply} reply/>) }
						</div>
						{ this.state.showCommentBox && 
							<textarea
								onChange={this.change}
								name="body"
								placeholder="Ok so basically...">
							</textarea> }
						<button className="inlineButton" onClick={this.submit}>
							<i class="fas fa-comment"></i> { !this.state.showCommentBox ? "Create reply" : "Send Reply" }
						</button>
						{	this.state.error.length > 0 && 
							<p className="formError">{ this.state.error }</p> }
					</div> }
			</div>
		);
	}
}