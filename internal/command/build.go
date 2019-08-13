package command

import "github.com/spf13/cobra"

const (
	HeaderFlagName     = "header"
	AWSSigV4FlagName   = "aws-sigv4"
	AWSProfileFlagName = "aws-profile"
	AWSRegionFlagName  = "aws-region"
	JSONBodyFlagName   = "json"
	OutputFileFlagName = "output-file"
)

// Build the root command for httpreq.
func Build() *cobra.Command {
	get := buildGet()
	post := buildPost()
	delete := buildDelete()

	root := &cobra.Command{
		Use:   "httpreq",
		Short: "httpreq <method> <route> [options]",
		Long: `Execute an HTTP request. Supported HTTP methods are GET, POST and DELETE.

Routes can have any of the following formats:
  * http[s]://host[:port]/path 	(use as is)
  * :port/path 			(assume http://localhost:port/path)
  * /path			(assume http://localhost:80/path

Headers are specified as a comma separated list of keypairs: --header name1(:|=)value1,name2(:|=)value2 ...
or specified multiple times: --header name1(:|=)value1 --header name2(:|=)value2`,
	}

	root.AddCommand(get, post, delete)
	return root
}

func buildGet() *cobra.Command {
	get := &cobra.Command{
		Use:   "get <url>",
		Short: "HTTP GET request.",
		Args:  cobra.ExactArgs(1),
		Run:   handleGet,
	}

	addCommonFlags(get)
	return get
}

func buildPost() *cobra.Command {
	post := &cobra.Command{
		Use:   `post <url> --json <body>`,
		Short: "HTTP POST request with a JSON body.",
		Long: `Make an HTTP POST request to the URL with a JSON body.
This command requires the --json flag, which should be a string conforming to valid JSON.`,
		Args: cobra.ExactArgs(1),
		Run:  handlePost,
	}

	post.Flags().String("json", "", "JSON body to use")
	addCommonFlags(post)
	return post
}

func buildDelete() *cobra.Command {
	delete := &cobra.Command{
		Use:   `delete <url>`,
		Short: "HTTP DELETE request.",
		Args:  cobra.ExactArgs(1),
		Run:   handleDelete,
	}

	addCommonFlags(delete)
	return delete
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice(
		HeaderFlagName,
		[]string{},
		`HTTP header to use in the request.
Value should be a keypair separated by equal sign (=) or colon (:), e.q. key=value.`)

	cmd.Flags().String(OutputFileFlagName, "", "Output the response body to the filename.")

	// AWS signature V4 flags
	cmd.Flags().Bool(AWSSigV4FlagName, false, "Use AWS signature V4 as authentication in the request. Requires the --aws-region option.")
	cmd.Flags().String(AWSRegionFlagName, "", "The AWS region to use in the AWS signature.")
	cmd.Flags().String(AWSProfileFlagName, "", "The name of an AWS profile in your AWS configuration. If not specified, environment variables are used.")
}
