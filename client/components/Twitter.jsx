import React, { Component } from "react";

export default class Twitter extends Component {
	constructor(props){
		super(props);
	}

	render(){
		return (
			<a href={this.props.URL}><img src="/static/twitter.svg"></img></a>
		)
	}
}