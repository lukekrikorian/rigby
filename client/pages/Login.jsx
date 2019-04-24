import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import Header from "../components/Header";
import Center from "../components/Center";
import { setLoggedIn, getSelf } from "../misc/user";
import Checkbox from "../components/Checkbox";
import Highlight from "../components/Highlight";

class Login extends Component {
	constructor(props){
		super(props);
		this.state = {
			username: '',
			password: '',
			saveSession: false,
			error: '',
			redirect: false,
		};
		this.handleChange = this.handleChange.bind(this);
		this.submit = this.submit.bind(this);
	}

	handleChange(event){
		const { target } = event;
		if (target.type === "checkbox") { 
			this.setState({ [target.name]: target.checked });
		} else {
			this.setState({ [target.name]: target.value });
		}
	}

	submit(event){
		event.preventDefault();

		const body = {
			username: this.state.username, 
			password: this.state.password,
			saveSession: this.state.saveSession
		}

		fetch("/api/login", {
			method: "POST",
			credentials: "same-origin",
			body: JSON.stringify(body)
		}).then(response => {
			if (!response.ok) {
				response.text().then(error => this.setState({ error }));
				return
			}
			setLoggedIn(true);
			getSelf();
			this.setState({ redirect: true });
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
					<Highlight>Login</Highlight>
					<form>
						<input 
							type="text"
							onChange={this.handleChange}
							name="username"
							placeholder="Username"/>
						<br/>
						<input 
							type="password"
							onChange={this.handleChange}
							name="password"
							placeholder="Password"/>
						<br/>
						<Checkbox
							onChange={this.handleChange}
							label="Remember Me"
							name="saveSession"/>
						<br/>
						<input
							type="submit"
							value="Submit"
							onClick={this.submit}/>
						<p className="formError">{this.state.error}</p>
					</form>
				</Center>
			</div>
		)
	}
}

export default Login;