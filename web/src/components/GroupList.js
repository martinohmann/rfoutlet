import React from 'react';
import PropTypes from 'prop-types';
import { List, NoItemsListItem } from './List';
import GroupListItem from './GroupListItem';

export default function GroupList(props) {
  const { groups } = props;

  return (
    <List>
      {groups.map(group =>
        <GroupListItem key={group.id} {...group} />
      )}
      {groups.length === 0 ? (
        <NoItemsListItem
          primary="No groups configured."
          secondary="Check your rfoutlet config."
        />
      ) : ''}
    </List>
  );
}

GroupList.propTypes = {
  groups: PropTypes.array.isRequired,
};
