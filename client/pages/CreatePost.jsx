import React, { Component } from "react";
import { isLoggedIn } from "../misc/user";
import { Redirect } from "react-router-dom";
import Header from "../components/Header";
import Center from "../components/Center";
import Highlight from "../components/Highlight";
import { timingSafeEqual } from "crypto";

class CreatePost extends Component {
	constructor(props){
		super(props);
		this.state = {
			title: '',
			body: '',
			error: '',
			redirect: '',
			placeholder: 'Markdown is supported. Rigby is designed for long-form content.',
			background: "cream"
		};
		this.handleChange = this.handleChange.bind(this);
		this.submit = this.submit.bind(this);
		this.handleDrag = this.handleDrag.bind(this);
		this.handleDragEnter = this.handleDragEnter.bind(this);
		this.handleDrop = this.handleDrop.bind(this);
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

	handleDrag(e){
		e.stopPropagation();
		e.preventDefault();
	}

	handleDragEnter(e){
		this.handleDrag(e);
		this.setState({ background: "off-cream" });
	}

	handleDrop(e){
		this.handleDrag(e);

		let reader = new FileReader();

		let file = e.dataTransfer.files[0];
		if (!file.type.startsWith("text")) { return; }

		let name = file.name.split(".").slice(0, -1).join(".");
		name = name.charAt(0).toUpperCase() + name.slice(1);

		this.setState({ title: name, background: "cream" });
		reader.onload = e => { this.setState({ body: e.target.result }); }
		
		reader.readAsText(file);
	}

	render(){
		if (this.state.redirect) { return <Redirect to={this.state.redirect}/> }
		if (!isLoggedIn()) { return <Redirect to="/login"/> }
		return (
			<div 
				onDrop={this.handleDrop}
				onDragEnter={this.handleDragEnter}
				onDragOver={this.handleDrag}
				style={{ background: `var(--${this.state.background})` }}>
				<Header/>
				<Center full-width>
					<Highlight>Create a post</Highlight>
					<form>
						<input 
							value={this.state.title}
							onChange={this.handleChange}
							type="text"
							placeholder="Title"
							name="title"/>
						<textarea
							value={this.state.body}
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