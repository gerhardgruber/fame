import * as React from 'react';
import { Layout } from 'antd';
import PageHeader from '../PageHeader';

abstract class Page<P = any> extends React.Component<P> {
  abstract pageTitle(): string;
  abstract renderContent(): JSX.Element;

  render(): JSX.Element {
    return (
      <Layout>
          <PageHeader name={this.pageTitle()} />
          <Layout style={{padding: '10px'}}>
            {this.renderContent()}
          </Layout>
      </Layout>
    )
  }
}

export default Page;