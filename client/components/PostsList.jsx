import React, { Component } from "react";
import PostCard from "./PostCard";

class PostsList extends Component {
	constructor(props){
		super(props);
	}

	render(){
		return (
			<div className="postsList">
				{ this.props.posts && this.props.posts.map(post => <PostCard post={post} showAuthor={this.props.showAuthor}/>) }
			</div>
		);
	}
}

export default PostsList;