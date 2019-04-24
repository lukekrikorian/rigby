import React, { Component } from "react";
import ReactMarkdown from "react-markdown";
import { Link } from "react-router-dom";
import Header from "../components/Header";
import Comment from "../components/Comment";
import Center from "../components/Center";
import { isLoggedIn } from "../misc/user";

class Post extends Component {
	constructor(props){
		super(props);
		this.state = {
			post: {},
			postID: props.match.params.id,
			commentBody: '',
			error: ''
		};
		this.changeComment = this.changeComment.bind(this);
		this.submitComment = this.submitComment.bind(this);
		this.likeButtonSubmit = this.likeButtonSubmit.bind(this);
	}

	likeButtonSubmit(event){
		event.preventDefault();
		fetch(`/api/vote/${this.state.post.ID}`, { 
			method: "POST",
			credentials: "same-origin"
		}).then(response => {
				if (response.ok) {
					this.setState({
						post: {
							...this.state.post,
							Votes: this.state.post.Votes + 1
						}
					});
				}
		}).catch(console.error)
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
		const { Body, Title, Author, Comments, Votes } = this.state.post;
		return (
			<div>
				<Header/>
				<Center>
					<div class="post">
						<h4 className="postAuthor">A work by <Link to={`/~${Author}`}>{ Author }</Link></h4>
						{ isLoggedIn() && 
							<h4 className="likeButton inlineButton" onClick={this.likeButtonSubmit}>
								<i className="fas fa-arrow-alt-circle-up fa-2x"></i> 
									{ (Votes || 0).toString().padStart(2, "0") } 
							</h4> }
						<h1 className="title">{ Title }</h1>
						<ReactMarkdown className="postBody">{ Body }</ReactMarkdown>
						{ Comments && 
							<div class="comments">
								{ Comments.map(comment => <Comment Comment={comment}/>) }
							</div> }
						{ isLoggedIn() &&
							<div style={{ marginBottom: 10, fontSize: "0.6em" }}>
								<textarea
									onChange={this.changeComment}
									name="commentBody"
									placeholder="Create a new comment"></textarea>
								<button
									className="inlineButton"
									onClick={this.submitComment}><i class="fas fa-comment"></i> Send</button>
								<p className="formError">{ this.state.error }</p>
							</div> }
					</div>
				</Center>
			</div> 
		)
	}
}

export default Post;