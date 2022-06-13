import React, {
  Component
} from "react";
import Header from "./components/Header";
import InputSearch from "./components/InputSearch"
import SearchEmoji from "./components/SearchEmoji"
import SearchResult from "./components/SearchResult"

class App extends Component {

  constructor(props) {
    super(props);
    this.state = {
      SearchEmoji: SearchEmoji("", 600)
    };
  }


  handleSearchChange = event =>{
    this.setState({
      SearchEmoji: SearchEmoji(event.target.value,20)
    });
  };

  render (){
    return (
      <div>
      <Header />
      <InputSearch textChange={this.handleSearchChange} />
      <SearchResult emojiData={this.state.SearchEmoji} />
      </div>
    );
  }
}

export default App;