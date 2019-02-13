import React, { Component } from "react";
import { Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import "../css/Comment.css";

export default class Comment extends Component {
	constructor(props){
		super(props);
	}

	render(){
		const comment = this.props.Comment;
		return (
			<div className={this.props.className || "comment"}>
				<p className="author"><Link to={`~${comment.Author}`}>{comment.Author}</Link> on {comment.Created.substring(0, 10)}</p>
				<ReactMarkdown allowedTypes={["paragraph", "text", "strong", "emphasis", "link"]} source={comment.Body}/>
				<div className="replies">
				{ comment.Replies &&
					comment.Replies.map(reply => <Comment className="reply" Comment={reply}/>) }
				</div>
			</div>
		);
	}
}