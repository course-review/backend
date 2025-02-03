package server

type CourseReviewService struct {
	db *sql.Queries
	pb.UnimplementedCourseReviewServer
}

func (s *server) GetUserData(ctx context.Context, in *pb.GetUserDataRequest) (*pb.GetUserDataResponse, error) {
	data, err := s.db.GetUserData(ctx, in.User)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserDataResponse{Data: data}, nil
}

func (s *server) GetReviews(ctx context.Context, in *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
    reviews, err := s.db.GetReviews(ctx, in.Course)
    if err != nil {
        return nil, err
    }
    return &pb.GetReviewsResponse{Reviews: reviews}, nil
}
