import React, { Component } from "react";
import "../css/TweetButton.css";

class TweetButton extends Component {
	constructor(props){
		super(props);
	}

	render(){
		return (
			<a target="_blank" href={`https://twitter.com/intent/tweet?url=${this.props.URL}`} className="tweetButton" rel="noopener">
				<img alt="Twitter Logo" src="/static/images/twitter.svg"/>
			</a>
		)
	}
}

export default TweetButton;