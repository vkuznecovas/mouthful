import { Link } from 'preact-router/match';
import Menu from '../Menu';
import style from './style';

const Header = () => (
  <header class={style.header}>
    <h1>Mouthful Admin Panel</h1>
    <Menu />
  </header>
);

export default Header;