import { h, Component } from 'preact';

export default class Button extends Component {
  render() {
    return (<div
      onClick={this.props.onClick} >
      {this.props.children}
    </div>);
  }
}