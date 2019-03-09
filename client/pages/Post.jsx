import React, { Component } from "react";
import ReactMarkdown from "react-markdown";
import { Link } from "react-router-dom";
import Header from "../components/Header";
import Comment from "../components/Comment";
import Center from "../components/Center";
import TweetButton from "../components/TweetButton";
import Heart from "../components/Heart";
import { isLoggedIn } from "../misc/user";
import "../css/Post.css";

class Post extends Component {
	constructor(props){
		super(props);
		this.state = {
			post: {},
			postID: props.match.params.id,
			showComment: false,
			commentBody: '',
			error: ''
		};
		this.toggleComment = this.toggleComment.bind(this);
		this.changeComment = this.changeComment.bind(this);
		this.submitComment = this.submitComment.bind(this);
	}
	
	toggleComment(event){
		event.preventDefault();
		this.setState({ showComment: !this.state.showComment });
	}

	changeComment(event){
		const { target } = event;
		this.setState({
			[target.name]: target.value
		});
	}

	submitComment(event){
		event.preventDefault();
		const body = {
			body: this.state.commentBody,
			postID: this.state.post.ID
		};
		fetch("/api/comments", {
			method: "POST",
			credentials: "same-origin",
			body: JSON.stringify(body)
		}).then(response => {
			if (!response.ok) {
				response.text().then(error => this.setState({ error }));
			} else { location.reload(); } 
		}).catch(console.error);
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
					{ isLoggedIn() && <Heart post={post.ID}/> }
					<ReactMarkdown source={post.Body} className={`${post.Comments && post.Comments.length > 0 && 'bottomBorder'} postBody`}/>
					<div class="comments">
						{ post.Comments && 
							post.Comments.map(comment => <Comment Comment={comment}/>) }
					</div>
					<button className="linkLike" onClick={this.toggleComment}>Create a comment</button>
					{ this.state.showComment &&
						<div>
							<textarea style={{ margin: 8, width: "24em", marginLeft: 0 }} onChange={this.changeComment} name="commentBody" placeholder="Ok so basically..."></textarea>
							<br/>
							<button style={{ marginRight: 4 }} className="linkLike" onClick={this.submitComment}>Send</button>
							<span>{this.state.error}</span>
					</div> }
				</Center>
			</div>
		)
	}
}

export default Post;