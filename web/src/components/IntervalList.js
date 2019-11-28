import React from 'react';
import PropTypes from 'prop-types';
import { List, NoItemsListItem } from './List';
import IntervalListItem from './IntervalListItem';

export default function IntervalList(props) {
  const { intervals, onDelete, onEdit, onToggle } = props;

  return (
    <List>
      {intervals.map(interval => (
        <IntervalListItem
          key={interval.id}
          interval={interval}
          onToggle={() => onToggle(interval)}
          onEdit={() => onEdit(interval)}
          onDelete={() => onDelete(interval)}
        />
      ))}
      {intervals.length === 0 ? (
        <NoItemsListItem
          primary="No intervals configured yet."
          secondary="Tap '+' to create one."
        />
      ) : ''}
    </List>
  );
}

IntervalList.propTypes = {
  intervals: PropTypes.array.isRequired,
  onDelete: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
  onToggle: PropTypes.func.isRequired,
};
