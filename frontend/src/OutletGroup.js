import React from 'react';
import { List } from 'material-ui/List';
import FlatButton from 'material-ui/FlatButton';
import { Toolbar, ToolbarGroup, ToolbarTitle } from 'material-ui/Toolbar';

import Outlet from './Outlet';
import makeApiRequest from './api.js'

class OutletGroup extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.handleOnButtonClick = this.handleOnButtonClick.bind(this);
    this.handleOffButtonClick = this.handleOffButtonClick.bind(this);
    this.handleToggleButtonClick = this.handleToggleButtonClick.bind(this);
    this.registerOutlet = this.registerOutlet.bind(this);

    this.styles = {
      toolbar: {
        paddingLeft: 16,
        paddingRight: 6,
      },
      title: {
        fontSize: 14,
      },
      button: {
        margin: 0,
        minWidth: 0,
      },
      label: {
        padding: 10,
        fontSize: 12,
      },
    };

    this.state = {
      outlets: [],
    }
  }

  registerOutlet(outlet) {
    var outlets = this.state.outlets;

    outlets.push(outlet);

    this.setState({ outlets: outlets });
  }

  updateOutletStates(action) {
    let data = { group_id: this.props.groupId };

    makeApiRequest(`/outlet_group/${action}`, data, result => {
      result.outlets.map((outlet, key) => {
        if (undefined !== this.state.outlets[key]) {
          this.state.outlets[key].setState({
            isEnabled: outlet.state === 1,
          });
        }

        return outlet;
      });
    });
  }

  handleOnButtonClick() {
    this.updateOutletStates('on');
  }

  handleOffButtonClick() {
    this.updateOutletStates('off');
  }

  handleToggleButtonClick() {
    this.updateOutletStates('toggle');
  }

  render() {
    return (
      <List>
        <Toolbar style={this.styles.toolbar}>
          <ToolbarGroup>
            <ToolbarTitle text={this.props.attributes.identifier} style={this.styles.title} />
          </ToolbarGroup>
          <ToolbarGroup>
            <FlatButton
              label="On"
              primary={true}
              style={this.styles.button}
              labelStyle={this.styles.label}
              onClick={this.handleOnButtonClick}
            />
            <FlatButton
              label="Off"
              secondary={true}
              style={this.styles.button}
              labelStyle={this.styles.label}
              onClick={this.handleOffButtonClick}
            />
            <FlatButton
              label="Toggle"
              style={this.styles.button}
              labelStyle={this.styles.label}
              onClick={this.handleToggleButtonClick}
            />
          </ToolbarGroup>
        </Toolbar>
        {this.props.attributes.outlets.map((attributes, key) =>
          <Outlet
            key={key}
            outletId={key}
            groupId={this.props.groupId}
            attributes={attributes}
            registerOutlet={this.registerOutlet}
          />
        )}
      </List>
    );
  }
}

export default OutletGroup;

