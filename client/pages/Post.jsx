import React, { Component } from "react";
import ReactMarkdown from "react-markdown";
import { Link } from "react-router-dom";
import Header from "../components/Header";
import Comment from "../components/Comment";
import Center from "../components/Center";
import TweetButton from "../components/TweetButton";
import Heart from "../components/Heart";
import "../css/Post.css";

export default class Post extends Component {
	constructor(props){
		super(props);
		this.state = {
			post: {},
			postID: props.match.params.id
		};
	}
	
	loadData() {
		fetch(`/api/posts/${this.state.postID}`)
			.then(response => response.json())
			.then(json => this.setState({ post: json }))
			.catch(console.error);
	}

	componentDidMount(){
		this.loadData();
	}

	render(){
		const { post } = this.state;
		return (
			<div>
				<Header/>
				<Center>
					<h1 className="title">{post.Title}</h1>
					<h4 className="postAuthor">By <Link to={`/~${post.Author}`}>{post.Author}</Link> on {post.Created && post.Created.substring(0, 10)}</h4>
					<TweetButton URL={`https://rigby.space/posts/${post.ID}`}/>
					<Heart post={post.ID}/>
					<ReactMarkdown source={post.Body} className={`${post.Comments && post.Comments.length > 0 && 'bottomBorder'} postBody`}/>
					<div class="comments">
						{ post.Comments && 
							post.Comments.map(comment => <Comment Comment={comment}/>) }
					</div>
				</Center>
			</div>
		)
	}
}