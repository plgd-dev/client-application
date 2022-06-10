import React from "react";
import PropTypes from "prop-types";

class InputSearch extends React.Component  {
    static propTypes = {
        textChange: PropTypes.func
    } 
    handleChange = event => {
        this.props.textChange(event);
    }
    render() {
        return (
            <div className="searchIconCircle material-icons" >
            <input type="text" onChange={this.handleChange} />
                <b className="searchIcon">search</b>
            </div>
        )
    }
}
export default InputSearch;