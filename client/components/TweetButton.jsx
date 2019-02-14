import React, { Component } from "react";
import "../css/TweetButton.css";

export default class TweetButton extends Component {
	constructor(props){
		super(props);
	}

	render(){
		return (
			<a target="_blank" href={`https://twitter.com/intent/tweet?url=${this.props.URL}`} className="tweetButton"><img src="/static/images/twitter.svg"></img></a>
		)
	}
}