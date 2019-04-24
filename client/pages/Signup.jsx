import React, { Component } from "react";
import Header from "../components/Header";
import Center from "../components/Center";
import Highlight from "../components/Highlight";
import { Redirect } from "react-router-dom";

class Signup extends Component {
	constructor(props){
		super(props);
		this.state = { 
			username: '', 
			password: '', 
			error: '', 
			redirect: false 
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
			username: this.state.username, 
			password: this.state.password
		}

		fetch("/api/signup", {
			method: "POST",
			credentials: "same-origin",
			body: JSON.stringify(body)
		}).then(response => {
			if (!response.ok) {
				response.text().then(error => this.setState({ error }));
			} else {
				fetch("/api/login", {
					method: "POST",
					credentials: "same-origin",
					body: JSON.stringify(body)
				}).then(() => { 
					window.isLoggedIn = true;
					this.setState({ redirect: true })
				})
			}
		}).catch(console.error);
	}

	render(){
		if (this.state.redirect) {
			return <Redirect to="/"/>
		}

		return (
			<div>
				<Header/>
				<Center>
					<Highlight>Sign Up</Highlight>
					<form>
						<input type="text" value={this.state.username} onChange={this.handleChange} name="username" placeholder="Username"/>
						<br/>
						<input type="password" value={this.state.password} onChange={this.handleChange} name="password" placeholder="Password"/>
						<br/>
						<input style={{ marginTop: 25 }} type="submit" value="Submit" onClick={this.submit}/>
						<p className={`formError`}>{this.state.error}</p>
					</form>
				</Center>
			</div>
		);

	}
}

export default Signup;