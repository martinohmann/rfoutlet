import React from 'react';
import { ListItem } from 'material-ui/List';
import Toggle from 'material-ui/Toggle';

import makeApiRequest from './api.js'

class Outlet extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.handleToggle = this.handleToggle.bind(this)

    this.state = {
      isEnabled: this.props.attributes.state === 1,
    };

    this.props.registerOutlet(this);
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      isEnabled: nextProps.attributes.state === 1,
    });
  }

  handleToggle(event, isEnabled) {
    let data = { group_id: this.props.groupId, outlet_id: this.props.outletId };

    makeApiRequest('/outlet/toggle', data, result => {
      this.setState({ isEnabled: result.state === 1 });
    });
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
        primaryText={this.props.attributes.identifier}
        rightToggle={toggle}
      />
    );
  }
}

export default Outlet;
