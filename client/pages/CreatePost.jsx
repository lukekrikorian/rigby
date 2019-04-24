import React, { Component } from "react";
import { isLoggedIn } from "../misc/user";
import { Redirect } from "react-router-dom";
import Header from "../components/Header";
import Center from "../components/Center";
import Highlight from "../components/Highlight";

class CreatePost extends Component {
	constructor(props){
		super(props);
		this.state = {
			title: '',
			body: '',
			error: '',
			redirect: '',
			placeholder: 'Markdown is supported. Rigby is designed for long-form content.'
		};
		this.handleChange = this.handleChange.bind(this);
		this.submit = this.submit.bind(this);
	}

	handleChange(event){
		const target = event.target;
		this.setState({ [target.name]: target.value });
	}

	submit(event){
		event.preventDefault();

		const body = {
			title: this.state.title,
			body: this.state.body
		};

		fetch("/api/post", {
			method: "POST", 
			credentials: "same-origin",
			body: JSON.stringify(body)
		}).then(response => {
			response.text().then(resp => {
				this.setState({
					[response.ok ? "redirect" : "error"]: resp
				});
			});
		}).catch(console.error);
	}

	render(){
		if (this.state.redirect) { return <Redirect to={this.state.redirect}/> }
		if (!isLoggedIn()) { return <Redirect to="/login"/> }
		return (
			<div>
				<Header/>
				<Center full-width>
					<Highlight>Create a post</Highlight>
					<form>
						<input 
							onChange={this.handleChange}
							type="text"
							placeholder="Title"
							name="title"/>
						<textarea
							style={{ height: 400 }}
							onChange={this.handleChange}
							placeholder={this.state.placeholder}
							name="body"/>
						<input
							type="submit"
							value="Submit"
							onClick={this.submit}/>
					</form>
					<p className="formError">{this.state.error}</p>
				</Center>
			</div>
		)
	}
}

export default CreatePost;