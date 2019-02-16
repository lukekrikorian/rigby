import React, { Component } from "react";
import { Link } from "react-router-dom";

class PostCard extends Component {
	constructor(props){
		super(props);
	}

	render(){
		const post = this.props.post;
		return(
			<div className="postCard" style={{ margin: 5 }}>
				↑{post.Votes} <Link to={`/posts/${post.ID}`}>{post.Title}</Link>
				{ this.props.showAuthor && 
					<span> By <Link to={`/~${post.Author}`}>{post.Author}</Link></span> }
			</div>
		)
	}
}

export default PostCard;