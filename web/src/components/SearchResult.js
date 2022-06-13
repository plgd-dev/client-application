import React, { Component } from "react";
import PropTypes from "prop-types";
import EmojiBox from "./EmojiBox";

class SearchResult extends Component {
	static propTypes = {
		emojiData: PropTypes.array,
	};

	render() {
		return (
			<div className="emojiResult">
				{this.props.emojiData.length ? (
					this.props.emojiData.map((data) => (
						<EmojiBox
							key={data.title}
							symbol={data.symbol}
							title={data.title}
						/>
					))
				) : (
					<h2>Not Found</h2>
				)}
			</div>
		);
	}
}

export default SearchResult;
