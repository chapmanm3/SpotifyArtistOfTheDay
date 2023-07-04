import axios from 'axios'

export interface ArtistInfoObject {
  ID: number,
  CreatedAt: string,
  UpdatedAt: string,
  DeletedAt: string,
  SpotifyUrl: string,
  SpotifyId: string,
  Image: string,
  Name: string,
  Uri: string,
}

interface ArtistInfoResponse {
  artist: ArtistInfoObject
}

export default function getTodaysArtist(): Promise<ArtistInfoResponse> {
  return axios.get('http://localhost:8080/api/artist/today', { withCredentials: true }).then((resp) => resp.data)
}
