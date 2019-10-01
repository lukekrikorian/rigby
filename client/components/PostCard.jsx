import React, { Component } from "react";
import { Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";

class PostCard extends Component {
	constructor(props){
		super(props);
	}

	render(){
		const post = this.props.post;
		const bodyArray = post.Body.replace("\n", "").split(" ");
		bodyArray.pop();
		return(
			<Link className="postCard" to={`/posts/${post.ID}`}>
				<h3>{ post.Title }</h3>
				<p>
					<ReactMarkdown
						allowedTypes={["text", "paragraph"]}
						unwrapDisallowed={true}>
							{ bodyArray.join(" ") + "..." }
					</ReactMarkdown>
				</p>
				<p className="postVotes">
					<i 
						style={{ verticalAlign: "middle" }} 
						className="arrowIcon fas fa-arrow-alt-circle-up fa-2x">
					</i>
					<span
						style={{ verticalAlign: "middle" }}>
							{ post.Votes !== undefined && post.Votes.toString().padStart(2, "0") }
					</span>
				</p>
			</Link>
		)
	}
}

export default PostCard;