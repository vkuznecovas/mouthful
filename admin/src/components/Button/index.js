const Button = ({ children, ...props }) => (
  <div {...props}>
    {children}
  </div>
);

export default Button;