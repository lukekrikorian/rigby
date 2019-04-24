import React, { Component } from "react";
import Header from "../components/Header";
import Center from "../components/Center";
import Comment from "../components/Comment";
import Highlight from "../components/Highlight";

class Conversation extends Component {
	constructor(props){
		super(props);
		this.state = {
			conversation: []
		};
	}

	loadData(){
		fetch("/api/conversation")
			.then(response => response.json())
			.then(json => this.setState({ conversation: json }))
			.catch(console.error);
	}

	componentDidMount(){ this.loadData(); }

	render(){
		return (
			<div>
				<Header/>
				<Center>
					<Highlight>Recent Conversations</Highlight>
					{ this.state.conversation && this.state.conversation.map(comment => <Comment Comment={comment} showPost/>) }
				</Center>
			</div>
		);
	}
}

export default Conversation;