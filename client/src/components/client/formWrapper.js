import { h, Component } from "preact";
import style from "./style";
import timeago from "./timeago"
import cookies from "./cookies"
import Form from "./form"



export default class FormWrapper extends Component {
    constructor(props) {
        super(props);
        this.getStyle = this.getStyle.bind(this);
    }
    getStyle(c) {
        return this.props.config.useDefaultStyle ? style[c] : c
    }
    render(props) {
        if (!this.props.comment.Confirmed) {
            return null
        }
        return <div>
            <input
                class={this.getStyle("mouthful_reply_button")}
                onClick={() => this.props.flipFormVisibility(this.props.comment.Id)}
                type="submit"
                value={this.props.visible ? "Close" : "Reply"}>
            </input>
            <Form id={this.props.comment.Id} config={this.props.config} visible={this.props.visible} author={this.props.author} comment={""} replyTo={this.props.replyTo} submitForm={this.props.submitForm} />
        </div>

    }

}
