import React, { Component } from "react";

class Checkbox extends Component {
	constructor(props){
		super(props);
		this.handleChange = this.handleChange.bind(this);
		this.state = {
			checked: this.props.checked === true,
		};
	}

	handleChange(e){
		e.target.name = this.props.name
		e.target.value = !this.state.checked;
		this.props.onChange(e);
		this.setState({ checked: !this.state.checked });
	}

	render(){
		let { checked } = this.state;
		let { label } = this.props;
		let checkedClass = `far fa-${this.state.checked ? "check-" : ""}square fa-2x`;
		return (
			<div>
				<input 
					className="hide"
					type="checkbox"
					checked={checked}>
				</input>
				<div 
					className="checkbox"
					onClick={this.handleChange}>
						<i className={checkedClass}></i>
						<span className="label">{ label }</span>
				</div>
			</div>
		)
	}

}

export default Checkbox;
