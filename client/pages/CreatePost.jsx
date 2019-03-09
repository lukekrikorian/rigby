import React, { Component } from "react";
import { isLoggedIn } from "../misc/user";
import { Redirect } from "react-router-dom";
import Header from "../components/Header";
import Center from "../components/Center";

import "../css/Form.css";

class CreatePost extends Component {
	constructor(props){
		super(props);
		this.state = {
			title: '',
			body: '',
			error: '',
			redirect: '',
			placeholder: 'Here we view the rare amphibian rigby, scientific name rigbaeus spacetus. The rigby is a nocturnal animal, only choosing to come out at night. Rigbys have only one natural predator, the loqual, scientific name postuslowqualitus. This creature aims to vore the rigby, from the inside out. Therefore, any loquals are to be shot on sight.'
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
				<Center>
					<h1>Create a post</h1>
					<form>
						<input onChange={this.handleChange} type="text" placeholder="Title" name="title"/>
						<br/>
						<textarea style={{ marginBottom: 6, minWidth: '24em', minHeight: '16em' }} onChange={this.handleChange} placeholder={this.state.placeholder} name="body"/>
						<br/>
						<input type="submit" value="Submit" onClick={this.submit}/>
					</form>
					<p>{this.state.error}</p>
				</Center>
			</div>
		)
	}
}

export default CreatePost;