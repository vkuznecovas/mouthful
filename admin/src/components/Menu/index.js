import { Link } from 'preact-router/match';

const Menu = () => (
  <nav>
    <Link href="/">All comments</Link>
    <Link href="/pending">Pending comments</Link>
    <Link href="/deleted">Deleted comments</Link>
  </nav>
);

export default Menu;