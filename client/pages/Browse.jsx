import React, { Component } from "react";
import Header from "../components/Header";
import Center from "../components/Center";
import PostsList from "../components/PostsList";

class Browse extends Component {
	constructor(props){ 
		super(props);
		this.state = {
			page: props.match.params.page,
			posts: [],
		};
	}

	loadData(){
		fetch(`/api/browse/${this.state.page}`)
			.then(response => response.json())
			.then(json => this.setState({ posts: json }))
			.catch(console.error);
	}

	componentDidMount(){ this.loadData(); }

	render(){
		return (
			<div>
				<Header/>
				<Center>
					<h1>Browsin'</h1>
					{ this.state.posts && <PostsList posts={this.state.posts}></PostsList> }
				</Center>
			</div>
		)
	}
}

export default Browse;