import { h, Component } from 'preact';
import { Link } from 'preact-router/match';

export default class Menu extends Component {
  render() {
    return (
      <nav>
        <Link href="/">All comments</Link>
        <Link href="/pending">Pending comments</Link>
        <Link href="/deleted">Deleted comments</Link>
      </nav>
    );
  }
}