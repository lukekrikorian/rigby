import React, { Component } from "react";
import { Link } from "react-router-dom";
import Header from "../components/Header";
import Center from "../components/Center";
import { isLoggedIn } from "../misc/user";

class Home extends Component {
	render () {
		if (isLoggedIn()) {
			return (
				<div>
					<Header/>
					<Center>
						<h1>Welcome back {window.user && window.user.Username}</h1>
						<Link to="/post">Create a post</Link>
						<br/>
						<Link to="/conversation">Recent conversation</Link>
						<br/>
						<Link to="/browse/recent">Recent posts</Link>
						<br/>
						<Link to="/browse/popular">Popular posts</Link>
						<br/>
						<Link to="/~me">Your profile</Link>
						<br/>
						<Link to="/logout">Logout</Link>
					</Center>
				</div>
			)
		}

		return (
			<div>
				<Header/>
				<Center>
					<h2 style={{ marginBottom: 0, color: "#332f2f" }}>Rigby</h2>
					<p style={{ marginTop: 5, marginBottom: 0 }}>a site where cool kids can discuss</p>
					<p style={{ marginTop: 4 }}><Link to="/signup">Sign up</Link> or <Link to="/login">Log in</Link> or <Link to="/browse/recent">Browse recent posts</Link></p>
				</Center>
			</div>
		)
	}
}

export default Home;