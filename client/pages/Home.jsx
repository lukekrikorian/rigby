import React, { Component } from "react";
import { Link } from "react-router-dom";
import Header from "../components/Header";
import Center from "../components/Center";
import { isLoggedIn } from "../misc/user";
import Highlight from "../components/Highlight";

class Home extends Component {
	render () {
		if (isLoggedIn()) {
			return (
				<div>
					<Header/>
					<Center>
						<Highlight>Welcome Back</Highlight>
						<div className="links">
							{ Object.entries({
								"/post": ["Create a Post", "pen"],
								"/conversation": ["Conversation", "comments"],
								"/browse/recent": ["Recent Posts", "clock"],
								"/browse/popular": ["Popular Posts", "arrow-alt-circle-up"],
								"/~me": ["Your Profile", "user-circle"],
								"/logout": ["Logout", "sign-out-alt"],
							}).map(([key, values]) => 
									<Link 
										className="homeLink inlineButton" 
										to={key}>
											<i class={`fas fa-${values[1]}`}></i> { values[0] }
									</Link>
							) }
						</div>
					</Center>
				</div>
			)
		}

		return (
			<div>
				<Header/>
				<Center>
					<Highlight>Rigby</Highlight>
					<p style={{ marginTop: 5, marginBottom: 0 }}>
							a site where cool kids can discuss
					</p>
					<p style={{ marginTop: 4 }}>
						<Link to="/signup">Sign up</Link> or <Link to="/login">Log in</Link> or <Link to="/browse/recent">Browse recent posts</Link>
					</p>
				</Center>
			</div>
		)
	}
}

export default Home;