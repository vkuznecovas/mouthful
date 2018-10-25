import { h, Component } from 'preact';

export default class Comment extends Component {
  constructor(props) {
    super(props);
    this.state = {
      a: null,
    };
  }

  render() {
    return (
        <div>{this.props.author} - {this.props.body}</div>
    );
  }
}