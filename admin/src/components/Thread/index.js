import { h, Component } from 'preact';

export default class Thread extends Component {
  constructor(props) {
    super(props);
    this.state = {
      a: null,
    };
  }

  render() {
    return (
      <div>
        <h1>THREAD</h1>
        {this.props.children}
      </div>
    );
  }
}