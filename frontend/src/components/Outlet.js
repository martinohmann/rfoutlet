import React from 'react';
import { ListItem } from 'material-ui/List';
import Toggle from 'material-ui/Toggle';

import { makeApiRequest, isOutletEnabled } from '../util'

class Outlet extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.handleToggle = this.handleToggle.bind(this)

    this.state = {
      isEnabled: isOutletEnabled(props.attributes),
    };

    this.props.registerOutlet(this);
  }

  componentWillReceiveProps(nextProps) {
    let outlet = nextProps.attributes;

    this.setState({ isEnabled: isOutletEnabled(outlet) });
  }

  handleToggle(event, isEnabled) {
    let data = { group_id: this.props.groupId, outlet_id: this.props.outletId };

    makeApiRequest('/outlet/toggle', data, outlet => {
      this.setState({ isEnabled: isOutletEnabled(outlet) });
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
