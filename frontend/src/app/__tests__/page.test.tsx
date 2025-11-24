import { render, screen } from '@testing-library/react';
import Home from '../page';

// Mock Next.js Image component
jest.mock('next/image', () => ({
  __esModule: true,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  default: ({ priority, quality, placeholder, loading, ...imgProps }: React.ComponentProps<'img'> & { priority?: boolean; quality?: number; placeholder?: string; loading?: string }) => {
    // Filter out Next.js-specific props that aren't valid HTML attributes
    // eslint-disable-next-line @next/next/no-img-element, jsx-a11y/alt-text
    return <img {...imgProps} />;
  },
}));

describe('Home Page', () => {
  it('renders the main heading', () => {
    render(<Home />);

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent(
      /To get started, edit the page.tsx file/i
    );
  });

  it('renders the Next.js logo', () => {
    render(<Home />);

    const logo = screen.getByAltText('Next.js logo');
    expect(logo).toBeInTheDocument();
  });

  it('renders the Vercel logo', () => {
    render(<Home />);

    const logo = screen.getByAltText('Vercel logomark');
    expect(logo).toBeInTheDocument();
  });

  it('contains links to documentation', () => {
    render(<Home />);

    const docsLink = screen.getByRole('link', { name: /documentation/i });
    expect(docsLink).toBeInTheDocument();
    expect(docsLink).toHaveAttribute('href', expect.stringContaining('nextjs.org/docs'));
  });

  it('contains link to templates', () => {
    render(<Home />);

    const templatesLink = screen.getByRole('link', { name: /templates/i });
    expect(templatesLink).toBeInTheDocument();
    expect(templatesLink).toHaveAttribute('href', expect.stringContaining('vercel.com/templates'));
  });

  it('contains link to learning center', () => {
    render(<Home />);

    const learningLink = screen.getByRole('link', { name: /learning/i });
    expect(learningLink).toBeInTheDocument();
    expect(learningLink).toHaveAttribute('href', expect.stringContaining('nextjs.org/learn'));
  });

  it('contains Deploy Now button', () => {
    render(<Home />);

    const deployButton = screen.getByRole('link', { name: /deploy now/i });
    expect(deployButton).toBeInTheDocument();
    expect(deployButton).toHaveAttribute('href', expect.stringContaining('vercel.com/new'));
  });

  it('has proper external link attributes', () => {
    render(<Home />);

    const externalLinks = screen.getAllByRole('link', { name: /deploy now|documentation/i });

    externalLinks.forEach(link => {
      expect(link).toHaveAttribute('target', '_blank');
      expect(link).toHaveAttribute('rel', 'noopener noreferrer');
    });
  });
});
