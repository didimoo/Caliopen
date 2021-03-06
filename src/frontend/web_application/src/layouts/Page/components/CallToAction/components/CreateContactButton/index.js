import { bindActionCreators, compose } from 'redux';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { withTranslator } from '@gandi/react-translate';
import Presenter from './presenter';

const mapDispatchToProps = dispatch => bindActionCreators({
  action: () => push('/new-contact'),
}, dispatch);

export default compose(
  connect(null, mapDispatchToProps),
  withTranslator()
)(Presenter);
