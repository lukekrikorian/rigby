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
		const pageName = (p => p === "recent" ? "Recent" : "Popular")(this.state.page);
		
		return (
			<div>
				<Header/>
				<Center>	
					<div class="highlight">
						<h1>{`${pageName} Posts`}</h1>
					</div>
					{ this.state.posts && <PostsList posts={this.state.posts}></PostsList> }
				</Center>
			</div>
		)
	}
}

export default Browse;