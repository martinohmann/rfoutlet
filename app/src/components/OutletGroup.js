import React from 'react';
import { List } from 'material-ui/List';
import IconButton from 'material-ui/IconButton';
import { Toolbar, ToolbarGroup, ToolbarTitle } from 'material-ui/Toolbar';

import Outlet from './Outlet';
import { styles } from '../config'
import { makeApiRequest, isOutletEnabled } from '../util'

class OutletGroup extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.registerOutlet = this.registerOutlet.bind(this);

    this.state = {
      outlets: [],
    }
  }

  registerOutlet(outlet) {
    var outlets = this.state.outlets;

    outlets.push(outlet);

    this.setState({ outlets: outlets });
  }

  handleButtonClick(action) {
    let data = {
      action: action,
      group_id: this.props.groupId
    };

    makeApiRequest(`/outlet_group`, data, result => {
      result.outlets.map((outlet, id) => {
        if (undefined !== this.state.outlets[id]) {
          this.state.outlets[id].setState({
            isEnabled: isOutletEnabled(outlet),
          });
        }

        return outlet;
      });
    });
  }

  render() {
    let { identifier, outlets } = this.props.attributes;
    let groupId = this.props.groupId;

    return (
      <List>
        <Toolbar style={styles.toolbar}>
          <ToolbarGroup>
            <ToolbarTitle text={identifier} style={styles.outletGroupTitle} />
          </ToolbarGroup>
          <ToolbarGroup>
            <IconButton
              iconClassName="material-icons"
              iconStyle={styles.buttonOff}
              onClick={(e) => this.handleButtonClick('off') }
            >
              flash_off
            </IconButton>
            <IconButton
              iconClassName="material-icons"
              iconStyle={styles.buttonOn}
              onClick={(e) => this.handleButtonClick('on') }
            >
              flash_on
            </IconButton>
            <IconButton
              iconClassName="material-icons"
              onClick={(e) => this.handleButtonClick('toggle') }
            >
              swap_horiz
            </IconButton>
          </ToolbarGroup>
        </Toolbar>
        {outlets.map((attributes, outletId) =>
          <Outlet
            key={outletId}
            outletId={outletId}
            groupId={groupId}
            attributes={attributes}
            registerOutlet={this.registerOutlet}
          />
        )}
      </List>
    );
  }
}

export default OutletGroup;

