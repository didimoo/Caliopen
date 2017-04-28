import { createSelector } from 'reselect';
import { bindActionCreators, compose } from 'redux';
import { connect } from 'react-redux';
import { withTranslator } from '@gandi/react-translate';
import { requestMessages, postActions } from '../../store/modules/message';
import { requestDiscussion } from '../../store/modules/discussion';
import Presenter from './presenter';

const messageByIdSelector = state => state.message.messagesById;
const discussionIdSelector = (state, ownProps) => ownProps.match.params.discussionId;

const mapStateToProps = createSelector(
  [messageByIdSelector, discussionIdSelector],
  (messagesById, discussionId) => ({
    discussionId,
    messages: Object.keys(messagesById)
      .map(messageId => messagesById[messageId])
      .filter(message => message.discussion_id === discussionId && message.is_draft !== true),
  })
);

const mapDispatchToProps = dispatch => bindActionCreators({
  requestDiscussion,
  requestMessages,
  setMessageRead: ({ message, isRead = true }) => {
    const action = isRead ? 'set_read' : 'set_unread';

    return postActions({ message, actions: [action] });
  },
}, dispatch);

export default compose(
  connect(mapStateToProps, mapDispatchToProps),
  withTranslator()
)(Presenter);
