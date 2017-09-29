import React from 'react';
import { ListItem } from 'material-ui/List';
import Toggle from 'material-ui/Toggle';

class Outlet extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.handleToggle = this.handleToggle.bind(this)

    this.state = {
      isEnabled: this.props.state === 1,
    };

    this.props.registerOutlet(this);
  }

  handleToggle(event, isEnabled) {
    this.setState({ isEnabled: isEnabled });
  }

  render() {
    var toggle = (
      <Toggle
        onToggle={this.handleToggle}
        toggled={this.state.isEnabled}
      />
    );

    return (
      <ListItem
        primaryText={this.props.identifier}
        rightToggle={toggle}
      />
    );
  }
}

export default Outlet;
