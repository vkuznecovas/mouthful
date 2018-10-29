import { h, Component } from 'preact';
import { Link } from 'preact-router/match';
import Menu from '../Menu';
import style from './style';

export default class Header extends Component {
  render() {
    return (
      <header class={style.header}>
        <h1>Mouthful Admin Panel</h1>
        <Menu />
      </header>
    );
  }
}