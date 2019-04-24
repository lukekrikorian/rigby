import React, { Component } from "react";
import Header from "../components/Header";
import Center from "../components/Center";
import PostsList from "../components/PostsList";

class Profile extends Component {
	constructor(props){
		super(props);
		this.state = {
			username: props.match.params.user,
			user: {}
		};
	}

	loadData(){
		fetch(`/api/users/${this.state.username}`)
			.then(response => response.json())
			.then(json => this.setState({ user: json }))
			.catch(console.error);
	}

	componentDidMount(){
		this.loadData();
	}

	render(){
		const { user } = this.state;
		return (
			<div>
				<Header/>
				<Center>
					<h1 class="profileHeader"><span style={{color: "var(--red)"}}>~</span>{user.Username}</h1>
					<h3 class="profileDate">Joined {user.Created && user.Created.substring(0, 10)}</h3>
					<PostsList posts={user.Posts} showAuthor={false}></PostsList>
				</Center>
			</div>
		)
	}
}

export default Profile;